package chatroom

import (
	"log"
	"net/http"
	"stbweb/core"
	"stbweb/lib/chatroom"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2/bson"

	"github.com/sirupsen/logrus"
)

//基本chat接口结构
type chat struct{}

//chatRoomBaseInfo 基本信息
type chatRoomBaseInfo struct {
	RoomName string //房间名称
	NumTotle int    //房间容量总人数
	RoomType string //房间类型
	Common   string //房间描述
}

func init() {
	core.RegisterFun("chat", new(chat), true)
}

//Post
func (ap *chat) Post(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionPost("create", new(chatRoomBaseInfo), createRoom) ||
		arge.APIInterceptionPost("condition", new(conditionWhere), selectCondition) {
		return
	}
}

//新建一个room对象，注意唯一号和房间对象都应该从对应池中获取，因为全局使用map保存，聊天websocket中也使用了map，考虑实际长度
func createRoom(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*chatRoomBaseInfo)
	if pm.NumTotle < 2 {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "房间人数不能少于两人"})
		return nil
	}
	roomID := chatroom.RoomIDPool.Get().(string)
	room := chatroom.RoomPool.Get().(chatroom.ChatRoom)

	ck := core.NewLock(pm.NumTotle, roomID)
	if !ck.GetLock(p.Usr) {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "获取进入房间资格失败"})
		return nil
	}

	room.RoomID = roomID
	room.HostName = p.Usr
	room.RoomName = pm.RoomName
	room.NumTotle = pm.NumTotle
	room.RoomType = pm.RoomType
	room.Common = pm.Common
	room.CreateTime = time.Now()
	core.RoomSets[roomID] = room

	if _, err := core.Mdb.InsertOne("chatroom", bson.M{"room_id": roomID, "host_name": p.Usr, "room_name": pm.RoomName, "room_type": pm.RoomType, "common": pm.Common}); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "room_id": roomID})
	return nil
}

func (*chat) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("leave", nil, userLeaveRoom) ||
		p.APIInterceptionGet("select", nil, selectRoom) {
		return
	}
}

//离开房间，注意如果是最后一个人，删除对应mongodb中的房间数据
//注意房主退出，直接清除所有成员
func userLeaveRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	if roomID == "" {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "roomID is nil"})
		return nil
	}
	if p.Usr == core.RoomSets[roomID].HostName {
		freedRoom(roomID)
	} else {
		core.RoomChatHub.Unregister(roomID, p.Usr)
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

//移交房主
func transferRoomHost(param interface{}, p *core.ElementHandleArgs) error {
	nextuser := p.Req.URL.Query().Get("nextuser")
	roomID := p.Req.URL.Query().Get("roomid")
	room := core.RoomSets[roomID]
	if p.Usr == room.HostName {
		room.HostName = nextuser
		if _, err := core.Mdb.UpdateOne("chatroom", bson.M{"roomID": roomID}, bson.M{"$set": bson.M{"host_name": nextuser}}); err != nil {
			logrus.WithFields(logrus.Fields{"mongo delete chat": err}).Error("freeRoom")
		}
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
		return nil
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false})
	return nil
}

//清理房间分三步，清理房间对象数据和对应的房间锁对象，以及断开房间内的连接
func clearRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	freedRoom(roomID)
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

func freedRoom(roomID string) {
	room := core.RoomSets[roomID]
	room.Clear()
	ck := core.RoomLocks[roomID]
	ck.Clear(roomID)
	core.RoomChatHub.UnregisterALL(roomID)
	//删除mongo房间
	if err := core.Mdb.DeleteDocument("chatroom", bson.M{"roomID": roomID}); err != nil {
		logrus.WithFields(logrus.Fields{"mongo delete chat": err}).Error("freeRoom")
	}
}

//
func selectRoom(param interface{}, p *core.ElementHandleArgs) error {
	page := p.Req.URL.Query().Get("page")
	t, err := strconv.Atoi(page)
	if err != nil {
		return err
	}
	skip := (t - 1) * 6
	res, err := core.Mdb.SelectAll("chatroom", bson.M{}, options.Find().SetSkip(int64(skip)))
	if err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": res})
	return nil
}

// type condition struct {
// 	Page      int
// 	Skip      int
// 	Limit     int
// 	OrderCol  string
// 	OrderType string

// 	Where interface{}
// }

// where中最多包含这些条件，用户对应表中字段，直接断言成bson.M
type conditionWhere struct {

	// RoomID     string
	// HostName   string
	RoomName string
	RoomType string
	Num      int
	// Common     string
	// CreateTime string

	Limit int
	Skip  int
	Sort  string
}

//这里直接使用options中条件方法来接受
func selectCondition(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*conditionWhere)
	res, err := core.Mdb.SelectAll("chatroom", setCondition(pm), setOption(pm))
	if err != nil {
		return err
	}
	log.Println(res)
	return nil
}

func setOption(wh *conditionWhere) *options.FindOptions {
	op := &options.FindOptions{}

	switch {
	case wh.Limit == 0:
		op.SetLimit(int64(wh.Limit))
		fallthrough
	case wh.Skip == 0:
		op.SetSkip(int64(wh.Skip))
		fallthrough
	case wh.Sort != "":
		op.SetSort(wh.Sort)
	}
	log.Println("op:", op)
	return op
}

func setCondition(wh *conditionWhere) bson.M {
	op := map[string]interface{}{}
	switch {
	case wh.Num != 0:
		op["Num"] = wh.Num
		fallthrough
	// case wh.RoomName == "":
	// 	log.Println("wh.RoomName:", wh.RoomName)
	// 	op["name"] = wh.RoomName
	// 	fallthrough
	case wh.RoomType != "":
		op["RoomType"] = wh.RoomType
	}
	log.Println("op:", op)
	return op
}
