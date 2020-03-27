package core

import (
	"net/http"
)

//Element 工作元素
type Element struct {
	Name    string
	Control *Controlle
}

//Handle 执行一个工作元素
//这里需要先判断是否有对应的controller，防止为空异常
func (e *Element) Handle(p *ElementHandleArgs) {
	if e.Control == nil {
		SendJSON(p.Res, http.StatusNotFound, nil)
		return
	}
	e.Control.Handle(p)
}

//ElementLoad 初始化element
func ElementLoad(elementName string) *Element {
	//这里应该使用yaml文件
	//从yaml文件中，使用elementName去对照取出对应所有该数据元素的对象内容
	//对象内容包括name，controllerName
	//然后name等基本属性就直接赋值给反馈的数据元素，获取的controllerName作为key，去全局controlles中找出本次请求对应的结构对象
	//这里缺少了yaml对照信息获取这一步
	return &Element{
		Name:    elementName,
		Control: controlles[elementName],
	}
}

//ElementHandle 处理一个http请求，确定一个element
//这里最后一个参数，对应的是元素名称，很重要，因为设计到路由内容和对应的方法，这里需要仔细考虑
func ElementHandle(w http.ResponseWriter, r *http.Request, elementName string) {
	//不能放下层进行判断，因为需要在认证检查之前返回

	ele := ElementLoad(elementName)
	arge := NewElementHandleArgs(w, r, ele)

	ele.Handle(arge)
}
