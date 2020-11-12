package chatroom

import (
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
	if pm.NumTotle < 2 {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": "房间人数不能少于两人"})
		return nil
	}
	roomID := roomIDPool.Get().(string)
	room := roomPool.Get().(chatRoom)

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

	roomSets[roomID] = room
	if err := room.save(); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "room_id": roomID})
	return nil
}

func (*chat) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("leave", nil, userLeaveRoom) {
		return
	}
}

//离开房间，注意如果是最后一个人，需要销毁对应nsq主题，删除对应mongodb中的房间数据
//注意房主退出，替换房主
func userLeaveRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	core.RoomChatHub.Unregister(roomID, p.Usr)
	return nil
}
