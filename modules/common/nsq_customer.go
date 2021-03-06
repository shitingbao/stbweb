//Package common 多个消费者在同一个主题消费时的几种情况
//消费情况
//1.在不同的通道，每个人收到的消息都是一样的，比如发了三次-1，2，3，每一个消费者都能收到三次，内容是1，2，3
//2.在相同的通道内，消费者是随机分配消息的，哪个先消化完一个消息，就继续执行，直到主题内没有消息为止
//开启顺序问题
//如果该主题还没有消费者，发送数据后不会丢失，知道有消费
//如果是上述这种情况，在多个消费者都需要完整的消息时，可能会出错
//因为当第一个连接该主题时，由于速度非常快，直接把所有消息都消化了，导致其他消费者刚连上该主题，已经没有数据了，导致群发消息只有单个消费者能接收到
//解决也很简单，先连接，再发送，不然就会丢失
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
	if _, err := snsq.NewNsqCustomer(tcpNsqdAddrr, "Insert", "channel3", &nsqHandler{nsqHandlerID: "three"}); err != nil {
		log.Println("ustomer3:", err)
		return
	}
}
