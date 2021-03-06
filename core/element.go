package core

import (
	"errors"
	"net/http"
	"stbweb/lib/rediser"
	"sync"
)

//Element 工作元素
type Element struct {
	Name        string
	Control     *Controlle
	elementPool *sync.Pool
}

//Handle 执行一个工作元素
//这里需要先判断是否有对应的controller，防止为空异常
func (e *Element) Handle(p *ElementHandleArgs) {
	if e.Control == nil {
		SendJSON(p.Res, http.StatusNotFound, "control is nil")
		return
	}
	e.Control.Handle(p)
}

//elementLoad 初始化element
func elementLoad(elementName string) *Element {
	//这里应该使用yaml文件
	//从yaml文件中，使用elementName去对照取出对应所有该数据元素的对象内容
	//对象内容包括name，controllerName
	//然后name等基本属性就直接赋值给反馈的数据元素，获取的controllerName作为key，去全局controlles中找出本次请求对应的结构对象
	//这里缺少了yaml对照信息获取这一步
	return &Element{
		Name:    elementName,
		Control: controlles[elementName],
		elementPool: &sync.Pool{
			New: func() interface{} {
				return &ElementHandleArgs{}
			},
		},
	}
}

//ElementHandle 处理一个http请求，确定一个element
//这里最后一个参数，对应的是元素名称，很重要，因为设计到路由内容和对应的方法，这里需要仔细考虑
func ElementHandle(w http.ResponseWriter, r *http.Request, elementName string) {
	usr, err := isExternalCall(elementName, r)
	if err != nil {
		SendJSON(w, http.StatusOK, SendMap{"msg": err.Error()})
		return
	}
	ele := elementLoad(elementName)
	arge := ele.elementPool.Get().(*ElementHandleArgs)
	arge.set(w, r, ele, usr)
	ele.Handle(arge)
	arge.clear()
	ele.elementPool.Put(arge)
}

//isExternalCall 判断该操作元素下的api是否可以外部调用
//这里还需要考虑到登录的用户长时间访问不需要登录的接口的情况
//接口中如果有token，应该都带上，这样就算访问非登录功能，也更新对应token状态，只不过不反馈是否成功更新标识
func isExternalCall(elementName string, r *http.Request) (string, error) {

	tokens := r.Header.Get("token")
	usr := rediser.GetUser(Rds, tokens)
	if controlleNames[elementName] { //判断该元素是否需要登陆后使用
		if tokens == "" {
			return "", errors.New("Refuse")
		}
		if usr == "" {
			return "", errors.New("token失效，请登录或者重新登录")
		}
		if err := rediser.MaintainActivity(Rds, tokens); err != nil {
			return "", errors.New("token失效，请登录或者重新登录")
		}
	} else {
		rediser.MaintainActivity(Rds, tokens)
	}
	return usr, nil
}
