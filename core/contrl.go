package core

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

//Controlle 保存所有定义的业务结构
var (
	controlles     = map[string]*Controlle{}
	controlleNames = map[string]bool{}
)

//Controlle 控制器结构
type Controlle struct {
	ControlleName string
	AppControlle  interface{}
}

//ElementHandleArgs http请求类型
type ElementHandleArgs struct {
	Req     *http.Request
	Res     http.ResponseWriter
	Red     *redis.Client
	Usr     string
	Element *Element
}

//BillGetEvent 工作元素
type BillGetEvent interface {
	Get(arge *ElementHandleArgs)
}

//BillPostEvent 工作元素
type BillPostEvent interface {
	Post(arge *ElementHandleArgs)
}

//isAPI 是否是api请求，是返回true
func (e *ElementHandleArgs) isAPI() bool {
	if len(e.Req.Header.Get(WebAPIHanderName)) > 0 {
		return true
	}
	return false
}

//NewElementHandleArgs 反馈一个工作元素类型
func NewElementHandleArgs(w http.ResponseWriter, r *http.Request, ele *Element, usr string) *ElementHandleArgs {
	return &ElementHandleArgs{
		Req:     r,
		Res:     w,
		Element: ele,
		Red:     Rds,
		Usr:     usr,
	}
}

//RegisterFun 注册一个功能,第二个参数为对应结构，应该使用new关键字新开辟对象，防止断言出错,第三个参数为是否是外部API，true为需要登录后使用
func RegisterFun(name string, ctr interface{}, isOut bool) {
	if name == "" || ctr == nil {
		LOG.Panic("app register err........")
	}
	register(&Controlle{
		ControlleName: name,
		AppControlle:  ctr,
	}, isOut)
}
func register(ctr *Controlle, isOut bool) {
	controlleNames[ctr.ControlleName] = isOut
	controlles[ctr.ControlleName] = ctr
}

//Handle 执行一个工作元素
//这里需要用到recover，因为如果业务类中只定义了get或者post其中一个，然后请求中地址对了，方法错了，这里就会异常,返回404，但是这里会输出panic
func (c *Controlle) Handle(arge *ElementHandleArgs) {

	switch arge.Req.Method {
	case "GET":
		defer func() {
			if err := recover(); err != nil {
				LOG.WithFields(logrus.Fields{"api-get": c.ControlleName}).Panic("api")
			}
		}()
		f, ok := c.AppControlle.(BillGetEvent)
		if !ok {
			LOG.Error("get BillGetEvent change error")
		}
		f.Get(arge)
	case "POST":
		defer func() {
			if err := recover(); err != nil {
				LOG.WithFields(logrus.Fields{"api-post": c.ControlleName}).Panic("api")
			}
		}()
		f, ok := c.AppControlle.(BillPostEvent)
		if !ok {
			LOG.Error("post BillPostEvent change error")
		}
		f.Post(arge)
	}
}

//URL中放入工作元素和一些检查的标识（_S这些安全检查），local host：//8088/input?_s=asgaoilo168hdhD4
//以工作元素名称做对照，找到对应实现get的对应结构对象
//将该对象断言为统一请求结构体内部（Element，就是都实现了get的interface）
//执行Element的get方法

//init中的工作，将所有实现了get的方法放入总controlle中注册，他是一个map，根据请求的工作元素用来对照查找
//使用时，新定义对应结构体，注册如全局controlle中，并实现对应方法（get，post）
