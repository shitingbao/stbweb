package chatroom

import (
	"encoding/json"
	"stbweb/lib/snsq"
	"sync"

	"github.com/nsqio/go-nsq"
)

//技术选型
//nsq消息队列，解决竞争
//mongodb频繁请求和松散的数据关系
//tcp连接，自由管理连接以及用户活动状态

//适当使用sync.pool重用chatRoom

//一个房间对象，包含基本信息和连接
//一个房间，对应一个nsq连接主题
type chatRoom struct {
	RoomID     string
	RoomName   string
	HostName   string
	NumTotle   int
	Num        int
	RoomType   string
	Common     string
	roomClient *nsq.Consumer
}

var (
	//房间池，新建房间应该从这里获取
	roomPool *sync.Pool
	//房间id，也从这里取，因为保存关系使用的是map，长度关系需要考虑
	roomIDPool *sync.Pool
	//tcp连接地址
	tcpNsqAddree = "127.0.0.1:4150"
)

func init() {
	roomPool = &sync.Pool{
		New: func() interface{} {
			return new(chatRoom)
		},
	}
}

//清理房间后，加入池（回收）,并且关闭nsq连接
func (c *chatRoom) clear() {
	c.RoomID = ""
	c.RoomName = ""
	c.HostName = ""
	c.NumTotle = 0
	c.Num = 0
	c.RoomType = ""
	c.Common = ""
	roomPool.Put(c)
	c.roomClient.Stop()
}

//nsq 消息handle
type nsqHandle struct{}

//nsq 消息传递结构
type nsqMes struct {
	ID   string
	Name string
}

func (s *nsqHandle) HandleMessage(mes nsq.Message) error {
	res := &nsqMes{}
	if err := json.Unmarshal(mes.Body, res); err != nil {
		return err
	}
	return nil
}

//只对应用户进入房间的逻辑操作，解决竞争关系使用队列
//删除房间时，清除该主题,主题使用唯一uuid编码，生成时机待定
//主题和通道统一使用一样的唯一号，方便对应，唯一号应该在生成房间时先生成
func newRoomClient(tc string) (*nsq.Consumer, error) {
	return snsq.NewNsqCustomer(tcpNsqAddree, tc, tc, &nsqHandle{})
}
