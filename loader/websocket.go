package loader

import (
	"encoding/json"
	"net/http"
	"stbweb/core"
	"stbweb/lib/rediser"
	"stbweb/lib/ws"
)

//initChatWebsocket 初始化websocket hub，开启消息处理循环
func initChatWebsocket() (chatHub, ctrlHub, cardHun *ws.Hub) {
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
	})

	cardHun = ws.NewHub(func(data []byte, hub *ws.Hub) error {
		msg := ws.Message{}
		if err := json.Unmarshal(data, &msg); err != nil {
			return err
		}
		hub.Broadcast <- msg
		return nil
	})
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

	return
}
