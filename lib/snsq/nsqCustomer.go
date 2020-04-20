package snsq

import (
	"errors"

	"github.com/nsqio/go-nsq"
)

//NewNsqCustomer 新建一个消费者,handle必须是实现了HandleMessage方法,内部连接，handle中接收数据
func NewNsqCustomer(tcpNsqdAddrr, topic, channel string, handle interface{}) error {
	con, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return err
	}
	// defer con.Stop()
	hd, ok := handle.(nsq.Handler)
	if !ok {
		return errors.New("handle type error")
	}
	con.AddHandler(hd)
	err = con.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		return err
	}
	return nil
}
