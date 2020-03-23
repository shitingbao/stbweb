package common

import (
	"net/http"
	"stbweb/core"
)

//AppExample 业务类
type AppExample struct{}

//accessPost 实际用来处理逻辑接收数据结构的类型
type accessPost struct {
	Name string
}

//localhost:3001/example
//header web-api : example
func init() {
	core.RegisterFun("example", new(AppExample)) //example 为url中匹配的工作元素名称
}

//Get 业务处理,get请求的例子
func (ap *AppExample) Get(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionGet("example", nil, appExamplef) { //example 为 header中web-api匹配的审核执行名称
		return
	}
}
func appExamplef(pa interface{}, content *core.ElementHandleArgs) error {
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": "this is example get"})
	return nil
}

//Post 业务处理,post请求的例子
func (ap *AppExample) Post(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionPost("example", new(accessPost), appPostExamplef) {
		return
	}
}
func appPostExamplef(pa interface{}, content *core.ElementHandleArgs) error {
	param := pa.(*accessPost) //这里使用指针断言来获取body内容，因为上面类型参数必须使用new关键字
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"post msg": param})
	return nil
}
