package loader

import (
	"encoding/json"
	"net/http"
	"stbweb/core"
	"stbweb/lib/game"
	"stbweb/lib/rediser"
	"stbweb/lib/ws"
	"strings"

	"github.com/sirupsen/logrus"
)

var roomChanPrefix = "room_chan_"

//initChatWebsocket 初始化websocket hub，开启消息处理循环
func initChatWebsocket() (chatHub, ctrlHub, cardHun *ws.Hub, roomChatHub *core.RoomChatHubSet) {
	chatHub = ws.NewHub(func(data []byte, hub *ws.Hub) error {
		msg := ws.Message{}
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		//原样消息发公告
		hub.Broadcast <- msg
		return nil
	})
	go chatHub.Run()
	//chat 聊天消息
	http.HandleFunc("/sockets/chat", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(rediser.GetUser(core.Rds, r.Header.Get("Sec-WebSocket-Protocol")), chatHub, w, r)
		game.RegisterAndStart(r.Header.Get("Sec-WebSocket-Protocol"), chatHub)
	})

	cardHun = ws.NewHub(game.ResponseOnMessage)
	go cardHun.Run()
	http.HandleFunc("/sockets/game", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(rediser.GetUser(core.Rds, r.Header.Get("Sec-WebSocket-Protocol")), chatHub, w, r)
	})

	ctrlHub = ws.NewHub(nil)
	//ctrl 控制消息
	go ctrlHub.Run()
	http.HandleFunc("/sockets/ctrl", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(rediser.GetUser(core.Rds, r.Header.Get("Sec-WebSocket-Protocol")), ctrlHub, w, r)
	})

	roomChatHub = core.NewChatHub(func(data []byte, hub *core.RoomChatHubSet) error {
		msg := core.ChatMessage{}
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		//原样消息发公告,注意这用BroadcastUser，只对该房间的人发送
		hub.BroadcastUser <- msg
		return nil
	})
	//ctrl 控制消息
	go roomChatHub.Run()
	http.HandleFunc("/room/chat", func(w http.ResponseWriter, r *http.Request) {
		//Sec中的值含义：token，房间号，用分号隔开,这里连接前需要获取锁，不能超过房间限定的人数
		info := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), "&")
		if len(info) != 2 {
			logrus.WithFields(logrus.Fields{"Sec-WebSocket-Protocol": "len should 2"}).Error("roomChatHub")
			return
		}
		token, roomID := info[0], info[1]
		user := rediser.GetUser(core.Rds, token)
		if user == "" {
			logrus.WithFields(logrus.Fields{"user": "nil"}).Error("roomChatHub")
			return
		}
		room := core.RoomSets[roomID]
		if room == nil || room.RoomID == "" {
			//id为空说明该房间不存在，或者已经被清理了
			logrus.WithFields(logrus.Fields{"Room": "nil"}).Error("roomChatHub")
			return
		}
		if room.HostName != user { //如果是房主本人连接不用获取锁，这是保证创建房间第一次房主连接，免得刚创建出来，位置都被其他人连接了
			ck, ok := core.RoomLocks[roomID]
			if !ok {
				return
			}
			if !ck.GetLock(user) { //获取锁
				return
			}
		}
		core.ServeChatWs(user, roomID, roomChatHub, w, r)
	})
	return
}
