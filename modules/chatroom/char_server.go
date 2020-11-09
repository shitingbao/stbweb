package chatroom

import (
	"encoding/json"
	"net"
	"stbweb/core"
	"stbweb/lib/snsq"

	"github.com/nsqio/go-nsq"
)

var (
//chat nsq连接地址
// chatNsqAddree = "127.0.0.1:4150"
)

func handleConnect(con net.Conn) {
	//处理消息逻辑
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
	return snsq.NewNsqCustomer(core.WebConfig.ChatNsqAddree, tc, tc, &nsqHandle{})
}
