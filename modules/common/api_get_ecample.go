//这个模块用来测试和例子展示

package common

import (
	"net/http"
	"stbweb/core"
	"stbweb/lib/excel"

	"github.com/Sirupsen/logrus"
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
	if arge.APIInterceptionGet("example", nil, appExamplef) || //example 为 header中web-api匹配的审核执行名称
		arge.APIInterceptionGet("excel", nil, excelExport) ||
		arge.APIInterceptionGet("excelparse", nil, excelparse) {
		return
	}
}
func excelExport(pa interface{}, content *core.ElementHandleArgs) error {
	rowData := []map[string]string{}
	da := make(map[string]string)
	da["one"] = "one"
	da["Two"] = "two"
	da["三"] = "三"
	da["4"] = "4"
	da["date"] = "1994-08-01"
	rowData = append(rowData, da)
	rowDatat := []map[string]string{}
	dc := make(map[string]string)
	dc["asdf"] = "asdf"
	dc["asdf"] = "asdf"
	dc["三"] = "三"
	dc["4"] = "4"
	dc["date"] = "1994-08-01"
	rowDatat = append(rowDatat, dc)
	if err := excel.CreateExcel("stb", rowData, rowDatat); err != nil {
		core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": err.Error()})
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": "this is excel get"})
	return nil
}

func excelparse(pa interface{}, content *core.ElementHandleArgs) error {
	res, err := excel.ExportParse("./assets/stb.xlsx", "Sheet2")
	if err != nil {
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, res)
	return nil
}
func appExamplef(pa interface{}, content *core.ElementHandleArgs) error {
	u := user{}
	if err := core.Ddb.QueryRow("SELECT name FROM user where name=?", "stb").Scan(&u.Name); err != nil {
		core.LOG.WithFields(logrus.Fields{"get user": err}).Error("user")
	}
	core.SendJSON(content.Res, http.StatusOK, core.SendMap{"msg": u.Name})
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
