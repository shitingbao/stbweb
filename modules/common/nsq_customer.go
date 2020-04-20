package common

import (
	"encoding/json"
	"fmt"
	"log"
	"stbweb/core"
	"stbweb/lib/snsq"

	"github.com/nsqio/go-nsq"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.0.1:4150"
)

func init() {
	core.RegisterFun("nsqcum", new(nsqHandler), false)
}

func (*nsqHandler) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("cumtomer", nil, nsqCustomer) {
		return
	}
}

//声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
type nsqHandler struct {
	//消息数
	msqCount int64
	//标识ID
	nsqHandlerID string
}

//实现HandleMessage方法
//message是接收到的消息
//这个函数一定要有，不需要手动调用
func (s *nsqHandler) HandleMessage(message *nsq.Message) error {
	//没收到一条消息+1
	s.msqCount++
	// //打印输出信息和ID
	// log.Println(s.msqCount, s.nsqHandlerID)
	//打印消息的一些基本信息
	result := &nsqMes{}
	if err := json.Unmarshal(message.Body, result); err != nil {
		panic(err)
	}
	//这里的时间指的是该消息进入队列的时间
	log.Println(s.msqCount, s.nsqHandlerID, "--data:", result)
	return nil
}

//消费函数
func nsqCustomer(param interface{}, p *core.ElementHandleArgs) error {
	go startCustomerChannel1() //开启第一个消息读取
	go startCustomerChannel2() //开启第二个消息读取
	go startCustomerChannel3() //开启第三个消息读取
	return nil
}

//原始方法
func startCustomerChannel1() {
	//初始化配置
	config := nsq.NewConfig()
	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("Insert", "channel1", config)
	if err != nil {
		fmt.Println(err)
	}
	//添加处理回调
	com.AddHandler(&nsqHandler{nsqHandlerID: "One"})
	//连接对应的nsqd
	err = com.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		fmt.Println(err)
	}
}

//原始方法二
func startCustomerChannel2() {
	//初始化配置
	config := nsq.NewConfig()
	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("Insert", "channel2", config)
	if err != nil {
		fmt.Println(err)
	}
	//添加处理回调
	com.AddHandler(&nsqHandler{nsqHandlerID: "two"})
	//连接对应的nsqd
	err = com.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		fmt.Println(err)
	}
}

//使用封装方法
func startCustomerChannel3() {
	if err := snsq.NewNsqCustomer(tcpNsqdAddrr, "Insert", "channel3", &nsqHandler{nsqHandlerID: "three"}); err != nil {
		log.Println("ustomer3:", err)
		return
	}
}
