package snsq

import (
	"encoding/json"
	"errors"

	"github.com/nsqio/go-nsq"
)

// 前提是安装后开启nsq服务
// nsqlookupd
// nsqd --lookupd-tcp-address=127.0.0.1:4160
// nsqadmin --lookupd-http-address=127.0.0.1:4161
//使用注意：开启数据发送前，应该确保消费者先连接，不然第一次可能会发生消息都被第一个消化完，导致第二个连接开启时获取不到消息
//使用同一个通道的消息消费者需要注意，那这个通道内的消息就是i大家公用的，一起在线就一起消费，如果没在线就会被其他消费者消费完
//创建消费者注意，创建连接后，有生产就会消费

//ProducerModel 一个生成者
type ProducerModel struct {
	Data  interface{}
	Topic string
	TPro  *nsq.Producer
}

//NewNsqProducerClient 返回一个生产者nsq连接，输入nsq连接地址
func NewNsqProducerClient(tcpNsqdAddrr string) (*ProducerModel, error) {
	config := nsq.NewConfig()
	tPro, err := nsq.NewProducer(tcpNsqdAddrr, config)
	if err != nil {
		return nil, err
	}
	return &ProducerModel{TPro: tPro}, nil
}

//Pulish 发送消息，data需要赋值在对象内部，主题和nsq连接生成对象，返回err，数据在内部json化,发送后清楚
func (np *ProducerModel) Pulish() error {
	defer func() {
		np.Data = nil
	}()
	if np.Topic == "" || np.Data == nil {
		return errors.New("主题和data数据不能为空")
	}
	da, err := json.Marshal(np.Data)
	if err != nil {
		return err
	}
	//发布消息
	return np.TPro.Publish(np.Topic, da)
}
