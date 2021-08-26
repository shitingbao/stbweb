package snsq

import (
	"errors"

	"github.com/nsqio/go-nsq"
)

//NewNsqCustomer 新建一个消费者,handle必须是实现了HandleMessage方法,内部连接，handle中接收数据
// 注意注意注意：handleMessage 中如果超过两分钟就会异常报错，导致消息重新回到队列重新执行
func NewNsqCustomer(tcpNsqdAddrr, topic, channel string, handle interface{}) (*nsq.Consumer, error) {
	con, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	// // defer con.Stop()
	hd, ok := handle.(nsq.Handler)
	if !ok {
		return nil, errors.New("handle type error")
	}
	con.AddHandler(hd)
	err = con.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
