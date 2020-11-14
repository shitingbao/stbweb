package chatroom

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/chatroom"
	"time"
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
	if arge.APIInterceptionPost("create", new(chatRoomBaseInfo), createRoom) {
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
	// if err := room.save(); err != nil {//mongdodb保存，待定
	// 	return err
	// }
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "room_id": roomID})
	return nil
}

func (*chat) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("leave", nil, userLeaveRoom) {
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
}
