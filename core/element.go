package core

import (
	"net/http"
)

//Controlle 保存所有定义的业务结构
var Controlle map[string]interface{}

//ElementHandleArgs http请求类型
type ElementHandleArgs struct {
	// Render   *render.Render
	Req *http.Request
	Res http.ResponseWriter
}

//Element 工作元素
type Element interface {
	Get(arge *ElementHandleArgs)
}

//NewElementHandleArgs 反馈一个工作元素类型
func NewElementHandleArgs(w http.ResponseWriter, r *http.Request) *ElementHandleArgs {
	return &ElementHandleArgs{
		Req: r,
		Res: w,
	}
}

//URL中放入工作元素和一些检查的标识（_S这些安全检查），local host：//8088/input?_s=asgaoilo168hdhD4
//以工作元素名称做对照，找到对应实现get的对应结构对象
//将该对象断言为统一请求结构体内部（Element，就是都实现了get的interface）
//执行Element的get方法

//init中的工作，将所有实现了get的方法放入总controlle中注册，他是一个map，根据请求的工作元素用来对照查找
//使用时，新定义对应结构体，注册如全局controlle中，并实现对应方法（get，post）
