package core

import (
	"net/http"
)

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
