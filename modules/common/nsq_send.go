package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

type nsqMes struct {
	Name string
	Age  int
	Num  int
}

// 主函数
func nsqSend() {
	//初始化配置
	config := nsq.NewConfig()
	tPro, err := nsq.NewProducer(tcpNsqdAddrr, config)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 100; i++ {
		//主题
		topic := "Insert"
		//主题内容
		// tCommand := strconv.Itoa(i)
		command := nsqMes{
			Name: "shitingbao",
			Age:  18,
			Num:  i,
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
}
