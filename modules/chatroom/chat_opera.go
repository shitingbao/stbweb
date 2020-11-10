package chatroom

import (
	"log"
	"net/http"
	"stbweb/core"
)

var (
	//romid对应一个room,保存所有的房间唯一号和房间对象的对应关系,主要是保存对应的nsq队列连接
	roomNsqClients map[string]chatRoom
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

	rClient, err := newRoomClient(roomID)
	if err != nil {
		return err
	}
	room.roomClient = rClient

	roomNsqClients[roomID] = room

	if err := room.save(); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "room_id": roomID})
	return nil
}

func (*chat) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("enter", nil, userEnterRoom) ||
		p.APIInterceptionGet("leave", nil, userLeaveRoom) {
		return
	}
}

//进入房间,竞争关系，使用nsq获取前几个值，反馈成功的标识和房间唯一号，如何在服务端将房间号和服务端收到的连接对应上
func userEnterRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	log.Println("id:", roomID)
	return nil
}

//离开房间，注意如果是最后一个人，需要销毁对应nsq主题，删除对应mongodb中的房间数据
func userLeaveRoom(param interface{}, p *core.ElementHandleArgs) error {
	roomID := p.Req.URL.Query().Get("roomid")
	log.Println("id:", roomID)
	return nil
}
