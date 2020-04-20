package common

import (
	"encoding/json"
	"fmt"
	"log"
	"stbweb/core"
	"stbweb/lib/snsq"
	"time"

	"github.com/nsqio/go-nsq"
)

func init() {
	core.RegisterFun("nsq", new(nsqMes), false)
}

func (*nsqMes) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("nsq", nil, nsqSend) ||
		p.APIInterceptionGet("stbnsq", nil, stSend) {
		return
	}
}

type nsqMes struct {
	Name string
	Age  int
	Num  int
}

var total = 0

//原始开启方法：开启数据发送前，应该确保消费者先连接，不然第一次可能会发生消息都被第一个消化完，导致第二个连接开启时获取不到消息
func nsqSend(param interface{}, p *core.ElementHandleArgs) error {
	//初始化配置
	config := nsq.NewConfig()
	tPro, err := nsq.NewProducer(tcpNsqdAddrr, config)
	if err != nil {
		fmt.Println(err)
	}
	for i := 0; i < 5; i++ {
		total++
		//主题
		topic := "Insert"
		//主题内容
		// tCommand := strconv.Itoa(i)
		command := nsqMes{
			Name: "shitingbao",
			Age:  18,
			Num:  total,
		}
		btData, err := json.Marshal(command)
		if err != nil {
			panic(err)
		}
		//发布消息
		err = tPro.Publish(topic, btData)
		if err != nil {
			fmt.Println(err)
		}
	}
	log.Println(time.Now().Format("2006-01-02 03:04:05"))
	return nil
}

//封装方法，注意事项同上
func stSend(param interface{}, p *core.ElementHandleArgs) error {
	pm, err := snsq.NewNsqProducerClient(tcpNsqdAddrr)
	if err != nil {
		return err
	}
	pm.Topic = "Insert"
	for i := 0; i < 5; i++ {
		total++
		pm.Data = nsqMes{
			Name: "shitingbao",
			Age:  18,
			Num:  total,
		}
		if err := pm.Pulish(); err != nil {
			return err
		}
	}
	return nil
}
