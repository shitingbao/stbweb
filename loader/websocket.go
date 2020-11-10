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
func initChatWebsocket() (chatHub, ctrlHub, cardHun *ws.Hub, roomChatHub *ws.ChatHub) {
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

	roomChatHub = ws.NewChatHub(nil)
	//ctrl 控制消息
	go roomChatHub.Run()
	http.HandleFunc("/room/chat", func(w http.ResponseWriter, r *http.Request) {
		// rediser.GetUser(core.Rds, r.Header.Get("Sec-WebSocket-Protocol"))
		//Sec中的值含义：用户名，房间号，标识，用分号隔开，需要检查标识的时效性，代表是否能进入房间，在chat模块中通过获取进入房间资格接口获取
		info := strings.Split(r.Header.Get("Sec-WebSocket-Protocol"), ":")
		if len(info) != 3 {
			logrus.WithFields(logrus.Fields{"Sec-WebSocket-Protocol": "len should 2"}).Error("roomChatHub")
			return
		}
		if u := core.Rds.Get(roomChanPrefix + info[0]).Val(); u == "" {
			logrus.WithFields(logrus.Fields{"连接资格无效": info[0], "房间号": info[1]}).Error("chat")
			return
		}
		ws.ServeChatWs(info[0], info[1], roomChatHub, w, r)
	})
	return
}
