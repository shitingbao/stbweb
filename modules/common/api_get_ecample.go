//这个模块用来测试和例子展示

package common

import (
	"log"
	"net/http"
	"stbweb/core"
	"stbweb/lib/task"
	"time"

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
	core.RegisterFun("example", new(AppExample), false) //example 为url中匹配的工作元素名称
}

//Get 业务处理,get请求的例子
func (ap *AppExample) Get(arge *core.ElementHandleArgs) {
	if arge.APIInterceptionGet("example", nil, appExamplef) ||
		arge.APIInterceptionGet("task", nil, taskExample) ||
		arge.APIInterceptionGet("sql", nil, sqlExample) { //example 为 header中web-api匹配的审核执行名称
		return
	}
}
func sqlExample(pa interface{}, content *core.ElementHandleArgs) error {
	stmt, err := core.Ddb.Prepare(`INSERT INTO task(
		task_id, 
		user,
		task_type,
		spec,
		Is_save,
		create_time,
		complete,
		execution_time) VALUES(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec("123", "123", "123", "123", true, time.Now(), true, time.Now()); err != nil {
		return err
	}
	core.SendJSON(content.Res, http.StatusOK, true)
	return nil
}
func taskExample(pa interface{}, content *core.ElementHandleArgs) error {
	log.Println("this is start task")
	ts := task.NewTask("sys", "clearMember", "0/2 * * * * ? ", func() error {
		log.Println("this is task==========")
		return nil
	})
	ts.Run()
	return nil
}

func appExamplef(pa interface{}, content *core.ElementHandleArgs) error {
	u := apiUser{}
	if err := core.Ddb.QueryRow("SELECT name FROM user where name=?", "stb").Scan(&u.Name); err != nil {
		logrus.WithFields(logrus.Fields{"get user": err}).Error("user")
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
