package chatroom

import (
	"encoding/json"
	"stbweb/core"
	"stbweb/lib/snsq"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

//nsq 消息handle
//缓存总人数和当前人数
//这里的并不代表真实的连接，因为前端拿到这个资格时，可能还没有连接过来，和实际的连接数目不一致也有可能，注意时效性
//注意获取资格后不使用的情况，这里的资格应该有时效性，不使用就失效
type nsqHandle struct {
	Total   int
	Current int
}

//nsq 消息传递结构
type nsqMes struct {
	User string
}

//这里对应的是多个生产者（多个连接，房间内的多个用户），一个消费者（一个房间对应一个Customer服务端主题）
func (s *nsqHandle) HandleMessage(mes nsq.Message) error {
	res := &nsqMes{}
	if err := json.Unmarshal(mes.Body, res); err != nil {
		return err
	}
	cn := userRoomChannel[res.User] //注意的是，如果该通道已经关闭，就是说对应请求已经超时反馈了，这里需要panic处理
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{"chan mes": "out time"}).Error("nsq")
		}
	}()
	if s.Current < s.Total {
		s.Current++
		cn <- true
	} else {
		cn <- false
	}
	return nil
}

//只对应用户进入房间的逻辑操作，解决竞争关系使用队列
//删除房间时，清除该主题,主题使用唯一uuid编码，生成时机待定
//主题和通道统一使用一样的唯一号，方便对应，唯一号应该在生成房间时先生成
//默认生成的房间人数为1
func newRoomClient(tc string, total int) (*nsq.Consumer, error) {
	return snsq.NewNsqCustomer(core.WebConfig.ChatNsqAddree, tc, tc, &nsqHandle{Total: total, Current: 1})
}
