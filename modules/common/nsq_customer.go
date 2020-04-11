package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nsqio/go-nsq"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.0.1:4150"
)

//声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
type NsqHandler struct {
	//消息数
	msqCount int64
	//标识ID
	nsqHandlerID string
}

//实现HandleMessage方法
//message是接收到的消息
//这个函数一定要有，不需要手动调用
func (s *NsqHandler) HandleMessage(message *nsq.Message) error {
	//没收到一条消息+1
	s.msqCount++
	//打印输出信息和ID
	log.Println(s.msqCount, s.nsqHandlerID)
	//打印消息的一些基本信息
	result := &nsqMes{}
	if err := json.Unmarshal(message.Body, result); err != nil {
		panic(err)
	}
	//这里的时间指的是该消息进入队列的时间
	log.Println("time:", time.Unix(0, message.Timestamp).Format("2006-01-02 03:04:05"), "--adree:", message.NSQDAddress, "--data:", result)
	return nil
}

//消费函数
func nsqCustomer() {
	//初始化配置
	config := nsq.NewConfig()
	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("Insert", "channel1", config)
	if err != nil {
		fmt.Println(err)
	}
	//添加处理回调
	com.AddHandler(&NsqHandler{nsqHandlerID: "One"})
	//连接对应的nsqd
	err = com.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Second * 25)
}
