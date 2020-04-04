//Package task 执行一个任务
//定时任务或者普通任务
//任务中所有处理的数据都需要有日志包处理，为了防止数据丢失而设计，可以反推实现数据恢复
//任务是异步的，相互不影响，并且相同的用户提交处理相同的任务
package task

import (
	"stbweb/core"

	"github.com/robfig/cron"
)

var job = cron.New()

//这里的任务应该产生日志，文件日志待定

const (
	//Export 导出
	Export = "et"
	//BatchEdit 批量修改
	BatchEdit = "bt"
	//UserClearList 用户缓存清理
	UserClearList = "ut"
)

//Task 任务对象
type Task struct {
	TaskID   string //任务Id
	User     string //任务所属user
	TaskType string //任务类型
	Spec     string //定时标记
	Func     func() //逻辑处理
}

func init() {
	job.Start()
}

//Stop job关闭
func Stop() {
	job.Stop()
}

//Run 运行一个task
func (t *Task) Run() {
	if core.Rds.HGet(t.User, t.TaskType).Val() != "" { //说明有相同的任务，不在执行
		return
	}
	core.Rds.HSet(t.User, t.TaskType, t.TaskID)

	job.AddFunc(t.Spec, t.Func)
}

//NewTaskRun 返回一个任务对象
func NewTaskRun(id, user, tasktype, spec string, function func()) *Task {
	return &Task{
		TaskID:   id,
		User:     user,
		TaskType: tasktype,
		Spec:     spec,
		Func:     function,
	}
}

// func task() {
// 	job.AddFunc("", func() {})
// 	// job.AddJob()

// }
