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
func ElementLoad(elementName string) *Element {
	//问题是这里的属性内容怎么获取
	//这里的这个elementName名称，和后面的请求的apiname没有关系，只是一个工作元素的名称，为后期他的属性内容做标识铺垫
	return &Element{
		Name:    elementName,
		Control: controlles[elementName],
	}
}

//ElementHandle 处理一个http请求，确定一个element
func ElementHandle(w http.ResponseWriter, r *http.Request, elementName string) {

	ele := ElementLoad(elementName)
	arge := NewElementHandleArgs(w, r, ele)

	ele.Handle(arge)
}
