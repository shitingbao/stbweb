package chatroom

import (
	"sync"

	"stbweb/core"

	"gopkg.in/mgo.v2/bson"

	"github.com/pborman/uuid"
)

//技术选型
//mongodb频繁请求和松散的数据关系，获取房间列表等
//自定义锁机制，类似java分段锁，一个房间内一把锁给一个用户使用，解决进入房间的竞争关系
//适当使用sync.pool重用RoomID和chatRoom，因为RoomID作为key在房间对象和锁key中都用到，并且使用了map保持关系，防止map没有释放，复用该值可适当减少并提升效率

//一个房间对象，包含基本信息
type chatRoom struct {
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
	roomPool.Put(c)

}

//保存房间后，加入mongodb
func (c *chatRoom) save() error {
	bm := bson.M{"roomID": c.RoomID,
		"RoomName": c.RoomName,
		"HostName": c.HostName,
		"NumTotle": c.NumTotle,
		"RoomType": c.RoomType,
		"Common":   c.Common,
	}
	if _, err := core.Mdb.InsertOne("room", bm); err != nil {
		return err
	}
	return nil
}

func (c *chatRoom) add() error {
	// res, err := core.Mdb.Selectone("room", bson.M{"roomID": c.RoomID})
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (c *chatRoom) done() {}

//删除一个房间，删除mongodb，删除nsq主题，清理chatRoom对象
func (c *chatRoom) delete() {

}
