package chatroom

import (
	"sync"
	"time"

	"github.com/pborman/uuid"
)

//技术选型
//mongodb频繁请求和松散的数据关系，获取房间列表等
//这里mongodb存储类型关系为，user，roomid，roomtype，对应用户在哪个房间，以及房间类型
//进入房间和离开房间只需要插入和删除就行，数据之间不要有关系，不需要update，免得引起不必要关联逻辑错误
//查询人数只需要查询roomid的count就行，不存实时在线总数，
//自定义锁机制，类似java分段锁，一个房间内一把锁给一个用户使用，解决进入房间的竞争关系
//适当使用sync.pool重用RoomID和chatRoom，因为RoomID作为key在房间对象和锁key中都用到，并且使用了map保持关系，防止map没有释放，复用该值可适当减少并提升效率

var (

	//RoomIDPool 房间id，也从这里取，因为保存关系使用的是map，长度关系需要考虑
	RoomIDPool *sync.Pool
)

func init() {

	RoomIDPool = &sync.Pool{
		New: func() interface{} {
			return uuid.NewUUID().String()
		},
	}

}

//ChatRoom 一个房间对象，包含基本信息
type ChatRoom struct {
	RoomID     string     //房间唯一id
	HostName   string     //房主名称，user
	RoomLock   sync.Mutex //主要是为了房间的销毁，该锁是每个房间的局部锁
	CreateTime time.Time
	RoomName   string //房间名称
	NumTotle   int    //房间容量总人数
	RoomType   string //房间类型
	Common     string //房间描述
}

//Clear 清理房间后,放回roomid以及删除对应mongo内的数据
//cf反调函数，为了等待mongo中删除完成，以及对应锁的释放
//如果在外部删除，可能出现的情况是，锁释放后，还没删除该roomid的数据，这时候该roomid复用并写入数据库，造成数据丢失，非线程安全
func (c *ChatRoom) Clear(cf func()) {
	c.RoomLock.Lock()
	defer c.RoomLock.Unlock()
	cf()
	RoomIDPool.Put(c.RoomID)
}

// //保存房间后，加入mongodb
// func (c *chatRoom) save() error {
// 	bm := bson.M{"roomID": c.RoomID,
// 		"RoomName": c.RoomName,
// 		"HostName": c.HostName,
// 		"NumTotle": c.NumTotle,
// 		"RoomType": c.RoomType,
// 		"Common":   c.Common,
// 	}
// 	if _, err := core.Mdb.InsertOne("room", bm); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *chatRoom) add() error {
// 	// res, err := core.Mdb.Selectone("room", bson.M{"roomID": c.RoomID})
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	return nil
// }

// func (c *chatRoom) done() {}

// //删除一个房间，删除mongodb，删除nsq主题，清理chatRoom对象
// func (c *chatRoom) delete() {

// }
