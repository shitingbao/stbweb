package nsq

import (
	"encoding/json"

	"github.com/nsqio/go-nsq"
)

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

//Pulish 发送消息，输入data，主题和nsq连接生成对象，返回err，数据在内部json化
func (np *ProducerModel) Pulish() error {
	da, err := json.Marshal(np.Data)
	if err != nil {
		return err
	}
	//发布消息
	return np.TPro.Publish(np.Topic, da)
}
