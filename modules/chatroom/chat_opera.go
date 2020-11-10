package chatroom

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"stbweb/core"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/pborman/uuid"
)

var (
	//romid对应一个room,保存所有的房间唯一号和房间对象的对应关系,主要是保存对应的nsq队列连接
	roomNsqClients map[string]chatRoom
	//用于每个用户从队列中获取反馈值，nsq和接口的通讯方式，接口中提交，队列中反馈是否成功获取进入房间的资格（true/false），最后用redis保持资格时效性
	userRoomChannel map[string]chan bool
	//redis中保持进入房间资格的前缀的key，后面跟上对应user
	roomChanPrefix = "room_chan_"
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

	rClient, err := newRoomClient(roomID, pm.NumTotle)
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
	if p.APIInterceptionGet("enter", nil, userEnterRoomQualification) ||
		p.APIInterceptionGet("leave", nil, userLeaveRoom) {
		return
	}
}

//获取进入房间的资格，反馈一个资格编号，存入redis，具有时效性,相互之间是竞争关系，使用nsq队列判断对应roomid的房间是否已满
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
	//判断是否该用户已经有资格连接，重置该连接即可
	if u := core.Rds.Get(roomChanPrefix + p.Usr).Val(); u != "" {
		core.Rds.SetNX(roomChanPrefix+p.Usr, u, time.Second)
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"isEnter_repeat": true, "qual": u})
		return nil
	}
	cn := make(chan bool, 1)
	userRoomChannel[p.Usr] = cn //将通道放入map，给队列服务端使用
	defer close(cn)

	room := roomNsqClients[roomID]
	//第一次判断(判断实际的连接)，nsq中进行第二次判断（缓存中判断），排除并发问题
	if core.RoomChatHub.RoomUserNum(roomID) > room.NumTotle {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "room is full", "isEnter": false})
		return nil
	}

	config := nsq.NewConfig()
	tPro, err := nsq.NewProducer(core.WebConfig.ChatNsqAddree, config)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": err, "isEnter": false})
		return err
	}

	da, err := json.Marshal(nsqMes{
		User: p.Usr,
	})
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"isEnter": false})
		return err
	}
	//发布消息
	tPro.Publish(roomID, da)
	ticker := time.NewTicker(time.Second) //这里一定要加入超时，防止死锁
	select {
	case isEnter := <-cn:
		if isEnter {
			qualification := uuid.NewUUID().String()
			core.Rds.SetNX(roomChanPrefix+p.Usr, qualification, time.Second) //设置过期时间，websocket使用时应该先判断该值是否过期
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"isEnter": true, "qual": qualification})
		} else {
			core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "enter room fail", "isEnter": false})
		}
	case <-ticker.C:
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "enter room outtime", "isEnter": false})
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
