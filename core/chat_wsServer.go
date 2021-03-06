package core

//改进以前发送某个特定消息的逻辑，需要轮询查找所有连接
//重写websocket，为了将连接分类，根据roomid，将每个房间的连接根据房间号来分，便于直接找到对应房间内的连接
//注意，由于使用房间号作为key分类连接，在使用全局消息时，该用户的状态应该时在房间外，所以roomid是空的，注意消息通道Broadcast的使用
import (
	"encoding/json"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Minute

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	pongWaitSet = time.Now()
)

//ChatMessage 管道中的消息
//user中为空"",则为全体发送，写入username则为指定发送，包括自己的信息
type ChatMessage struct {
	User     string
	RoomID   string
	Data     interface{}
	DateTime time.Time
}

//NewChatHub 分配一个新的Hub，使用前先获取这个hub对象
func NewChatHub(onEvent OnMessageFuncChat) *RoomChatHubSet {
	return &RoomChatHubSet{
		Broadcast:     make(chan ChatMessage),           //包含要想向前台传递的数据，内部使用chan通道传输
		BroadcastUser: make(chan ChatMessage),           //包含对应指定房间的消息
		register:      make(chan *ChatClient),           //有新的连接，将放入这里
		unregister:    make(chan *ChatClient),           //断开连接加入这
		clients:       make(map[string]([]*ChatClient)), //包含所有的client连接信息
		OnMessage:     onEvent,
	}
}

//OnMessageFuncChat 接收到消息触发的事件
type OnMessageFuncChat func(message []byte, hub *RoomChatHubSet) error

// ChatClient is a middleman between the websocket connection and the hub.
//可增加一个用户属性，用来区分不同的连接，便于在发送的时候区分发送，不走同一个频道，这样就可以分为全局频道和局部频道
//需要登录配合，可以将用户登陆时保存在cookie中，在注册client时获取
type ChatClient struct {
	hub    *RoomChatHubSet
	user   string
	roomID string
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//服务端主动退出的标识
	//由于write和read只能等待超时后释放，当服务端主动关闭socket连接时，并不能及时清理conn
	//该信号在run中的unregister中放入，注册时保留一位缓存
	outSign chan bool
}

//RoomChatHubSet maintains the set of active clients and broadcasts messages to the
//clients.
type RoomChatHubSet struct {
	// Registered clients.
	clients map[string]([]*ChatClient)

	//Broadcast 公告消息队列
	Broadcast chan ChatMessage
	//用户私人消息队列
	BroadcastUser chan ChatMessage
	//OnMessage 当收到任意一个客户端发送到消息时触发
	OnMessage OnMessageFuncChat

	// Register requests from the clients.
	register chan *ChatClient

	// Unregister requests from clients.
	unregister chan *ChatClient
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *ChatClient) writePump() {
	ticker := time.NewTicker(pingPeriod) //设置定时
	defer func() {
		ticker.Stop()
		c.conn.Close()
		if lk := RoomLocks[c.roomID]; lk != nil { //可能在退出的时候，整个room对象都被清理了，就不用去放回了，业务情况是房主退出
			lk.FreedLock(c.user) //退出并释放锁，这个函数在读或者写中执行一次即可，不然就会每次端断开都有两次信号放回
		}
		if err := Mdb.DeleteDocument("chat", bson.M{"user": c.user}); err != nil {
			logrus.WithFields(logrus.Fields{"mongo delete chat": err}).Error("websocket")
		}
	}()
	for {
		select {
		case message, ok := <-c.send: //这里send里面的值时run里面获取的，在这里才开始实际向前台传值
			//这里add时间，逻辑上是因为发言越多，可持续超时时间就越长，一句话都不说，那就5分钟就滚犊子了
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //设置写入的死亡时间，相当于http的超时
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{}) //如果取值出错，关闭连接，设置写入状态，和对应的数据
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage) //以io的形式写入数据，参数为数据类型
			if err != nil {
				logrus.WithFields(logrus.Fields{"sock write:": err}).Error("websockets")
				return
			}
			w.Write(message) //写入数据，这个函数才是真正的想前台传递数据

			if err := w.Close(); err != nil { //关闭写入流
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait)) //心跳包，下面ping出错就会报错退出，断开这个连接
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-c.outSign:
			return
		}
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *ChatClient) readPump() {
	defer func() {
		c.hub.unregister <- c //读取完毕后注销该client
		c.conn.Close()
		// RoomLocks[c.roomID].FreedLock(c.user) //退出并释放锁,这个函数在读或者写中执行一次即可，不然就会每次端断开都有两次
	}()
	c.conn.SetReadLimit(maxMessageSize)
	pongTime := time.Now().Add(pongWait)
	pongWaitSet = pongTime           //设置最大读取容量
	c.conn.SetReadDeadline(pongTime) //设置读取死亡时间
	c.conn.SetPongHandler(func(string) error {
		t := time.Now().Add(time.Now().Sub(pongWaitSet))
		c.conn.SetReadDeadline(t)
		return nil
	}) //响应事件的设置，收到响应后，重新设置死亡时间
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logrus.WithFields(logrus.Fields{"ReadMessageError": err}).Info("websocket")
			}
			break
		}
		if c.hub.OnMessage != nil { //执行回调函数
			if err := c.hub.OnMessage(message, c.hub); err != nil {
				logrus.WithFields(logrus.Fields{"ReadOnMessageError": err}).Error("websocket")
				break
			}
		}
	}
}

//Run 开始消息读写队列，无限循环，应该用go func的方式调用
func (h *RoomChatHubSet) Run() {
	if err := recover(); err != nil {
		logrus.WithFields(logrus.Fields{"run": err}).Error("RoomChatHubSet")
	}
	for {
		select {
		case client := <-h.register: //客户端有新的连接就加入一个
			h.clients[client.roomID] = append(h.clients[client.roomID], client)
			// logrus.Info("当前连接增加，", "总连接数为：", len(h.clients))
		case client := <-h.unregister: //客户端断开连接，client会进入unregister中，直接在这里获取，删除一个
			ct := h.clients[client.roomID]
			for i, v := range ct {
				if v.user == client.user {
					close(client.send) //关闭对应连接
					h.clients[client.roomID] = append(h.clients[client.roomID][:i], h.clients[client.roomID][i+1:]...)
				}
			}
		case mes := <-h.Broadcast: //将数据发给所有连接中的send，用来发送全局消息，如系统提示消息或者全世界喊话
			data, err := json.Marshal(mes)
			if err == nil {
				for _, cts := range h.clients {
					for _, vclient := range cts {
						vclient.send <- data
					}
				}
			} else {
				logrus.WithFields(logrus.Fields{"json Marshal": err}).Error("websocket")
			}
		case mes := <-h.BroadcastUser: //将数据发给连接中的send，用来发送
			data, err := json.Marshal(mes)
			if err == nil {
				ct := h.clients[mes.RoomID]
				for _, client := range ct { //clients中保存了所有的客户端连接，循环所有连接给与要发送的数据
					client.send <- data
				}
			} else {
				logrus.WithFields(logrus.Fields{"json": err}).Error("BroadcastUser")
			}
		}
	}
}

//Len 返回房间数量
func (h *RoomChatHubSet) Len() int {
	return len(h.clients)
}

//RoomUserNum 反馈房间内人数
func (h *RoomChatHubSet) RoomUserNum(roomid string) int {
	return len(h.clients[roomid])
}

//Unregister 主动退出一个连接
func (h *RoomChatHubSet) Unregister(roomid, user string) {
	for i, v := range h.clients[roomid] {
		if v.user == user {
			h.clients[roomid][i].outSign <- true
		}
	}
}

//UnregisterALL 主动退出一个房间内的所有连接
//做这一步不用当心过程中有新连接，room对象中使用了锁，并在连接前后判断，在清理前连接，这里直接清除，清理后连接，逻辑中判断room内容为空，不将连接加入，直接return，行260
func (h *RoomChatHubSet) UnregisterALL(roomid string) {
	for i := range h.clients[roomid] {
		h.clients[roomid][i].outSign <- true
	}
}

// ServeChatWs handles websocket requests from the peer.
func ServeChatWs(user, roomID string, hub *RoomChatHubSet, w http.ResponseWriter, r *http.Request) {
	room := RoomSets[roomID]
	if room == nil {
		return
	}
	room.RoomLock.Lock()
	if room.HostName == "" || room.RoomID == "" {
		//房主为空说明该房间已经销毁了，id为空说明该房间不存在，或者已经被清理了
		room.RoomLock.Unlock()
		logrus.WithFields(logrus.Fields{"room": "HostName or RoomID is nil"}).Error("websocket")
		return
	}
	h := http.Header{}
	pro := r.Header.Get("Sec-WebSocket-Protocol")
	h.Add("Sec-WebSocket-Protocol", pro)   //带有websocket的Protocol子header需要传入对应header，不然会有1006错误
	conn, err := upgrader.Upgrade(w, r, h) //返回一个websocket连接
	if err != nil {
		logrus.WithFields(logrus.Fields{"connect": err}).Info("websocket")
		room.RoomLock.Unlock()
		return
	}

	//生成一个client，里面包含用户信息连接信息等信息
	client := &ChatClient{hub: hub, user: user, roomID: roomID, conn: conn, send: make(chan []byte, 256), outSign: make(chan bool, 1)}
	client.hub.register <- client //将这个连接放入注册，在run中会加一个
	room.RoomLock.Unlock()
	if _, err := Mdb.InsertOne("chat", bson.M{"roomId": room.RoomID, "user": user}); err != nil {
		logrus.WithFields(logrus.Fields{"mongo": err}).Error("websocket")
	}
	//注意这俩个读写方法都是持续监听的，一定要放最后
	go client.writePump() //新开一个写入，因为有一个用户连接就新开一个，相互不影响，在内部实现心跳包检测连接，详细看函数内部
	client.readPump()     //读取websocket中的信息，详细看函数内部
	//这一步解锁注意，上面已经有判断

}
