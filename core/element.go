package core

import "net/http"

//Element 工作元素
type Element struct {
	Name    string
	Control *Controlle
}

//Handle 执行一个工作元素
func (e *Element) Handle(arge *ElementHandleArgs) {
	e.Control.Handle(arge)

}

//ElementLoad 初始化element
func ElementLoad(elementName string) Element {
	//问题是这里的属性内容怎么获取

	return Element{
		Name:    elementName,
		Control: controlles[elementName],
	}
}

//ElementHandle 处理一个http请求，确定一个element
func ElementHandle(w http.ResponseWriter, r *http.Request, elementName string) {
	arge := NewElementHandleArgs(w, r)
	ele := ElementLoad(elementName)
	ele.Handle(arge)
}
