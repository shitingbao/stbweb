package chatroom

import (
	"errors"
	"log"
	"net/http"
	"stbweb/core"
)

var (
	//romid对应一个room,保存所有的房间唯一号和房间对象的对应关系
	roomSets map[string]chatRoom
)

//基本chat接口结构
type chat struct{}

func init() {
	core.RegisterFun("chat", new(chat), true)
}

//Post
func (ap *chat) Post(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionPost("create", new(chatRoomBaseInfo), createRoom) {
		return
	}
}

//新建一个room对象，注意唯一号和房间对象都应该从对应池中获取，因为全局使用map保存，聊天websocket中也使用了map，考虑实际长度
func createRoom(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*chatRoomBaseInfo)

	roomID := roomIDPool.Get().(string)
	room := roomPool.Get().(chatRoom)

	room.RoomID = roomID
	room.HostName = p.Usr
	room.RoomName = pm.RoomName
	room.NumTotle = pm.NumTotle
	room.RoomType = pm.RoomType
	room.Common = pm.Common

	roomSets[roomID] = room
	if err := room.save(); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "room_id": roomID})
	return nil
}

func (*chat) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("enter", nil, userEnterRoomQualification) ||
		p.APIInterceptionGet("leave", nil, userLeaveRoom) {
		return
	}
}

//获取进入房间的资格，反馈一个资格编号，存入redis，具有时效性,相互之间是竞争关系，使用nsq队列判断对应roomid的房间是否已满
//redis中的锁应该对应user
//一定要先判断usr是否为空，因为usr是作为消息传递的基础来的
//使用user对应的chan来接受反馈的数据
func userEnterRoomQualification(param interface{}, p *core.ElementHandleArgs) error {
	if p.Usr == "" {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "user is nil", "isEnter": false})
		return errors.New("user not nil")
	}

	roomID := p.Req.URL.Query().Get("roomid")
	if roomID == "" {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "roomid is nil", "isEnter": false})
		return errors.New("roomid not nil")
	}
	roomLock := core.RoomLocks[roomID]
	if roomLock.GetLock(p.Usr) { //true为获取成功
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	}
	return nil
}

//离开房间，注意如果是最后一个人，需要销毁对应nsq主题，删除对应mongodb中的房间数据
//注意房主退出，替换房主
func userLeaveRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	log.Println("id:", roomID)
	return nil
}
