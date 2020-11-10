package chatroom

import (
	"sync"

	"stbweb/core"

	"gopkg.in/mgo.v2/bson"

	"github.com/nsqio/go-nsq"
	"github.com/pborman/uuid"
)

//技术选型
//nsq消息队列，解决竞争
//mongodb频繁请求和松散的数据关系
//tcp连接，自由管理连接以及用户活动状态

//适当使用sync.pool重用chatRoom

//一个房间对象，包含基本信息和连接
//一个房间，对应一个nsq连接主题
type chatRoom struct {
	chatRoomBase
	// hub        hubClient
	roomClient *nsq.Consumer //每个房间对应的队列连接，销毁房间时断开
}

//保存所有成员，以及对应tcp连接,用户user对应自己的tcp连接对象，这个是客户端连接
// type hubClient map[string]*net.TCPConn

//mongodb只保存这部分，用于查询即可
type chatRoomBase struct {
	RoomID   string //房间唯一id
	HostName string //房主名称，user
	chatRoomBaseInfo
}

type chatRoomBaseInfo struct {
	RoomName string //房间名称
	NumTotle int    //房间容量总人数
	RoomType string //房间类型
	Common   string //房间描述
}

var (
	//房间池，新建房间应该从这里获取，放回前需要先调用clear
	roomPool *sync.Pool
	//房间id，也从这里取，因为保存关系使用的是map，长度关系需要考虑
	roomIDPool *sync.Pool
	//tcp连接复用，退出房间放回,分配一个连接，退出房间或者超时都应该断开

)

func init() {
	roomPool = &sync.Pool{
		New: func() interface{} {
			return new(chatRoom)
		},
	}
	roomIDPool = &sync.Pool{
		New: func() interface{} {
			return uuid.NewUUID().String()
		},
	}

}

//清理房间后，加入池（回收）,并且关闭nsq连接
func (c *chatRoom) clear() {
	c.RoomID = ""
	c.RoomName = ""
	c.HostName = ""
	c.NumTotle = 0
	c.RoomType = ""
	c.Common = ""
	// c.hub = make(hubClient)
	c.roomClient.Stop()
	roomPool.Put(c)

}

//保存房间后，加入mongodb
func (c *chatRoom) save() error {
	bm := bson.M{"roomID": c.RoomID,
		"RoomName": c.RoomName,
		"NumTotle": c.NumTotle,
		"RoomType": c.RoomType,
		"Common":   c.Common,
	}
	if _, err := core.Mdb.InsertOne("room", bm); err != nil {
		return err
	}
	return nil
}

//删除一个房间，删除mongodb，删除nsq主题，清理chatRoom对象
func (c *chatRoom) delete() {

}
