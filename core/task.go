package core

import (
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

//TaskRun 任务对象
type TaskRun struct {
	TaskID   string
	User     string
	TaskType string
	Spec     string
	Func     func()
}

func init() {
	job.Start()
}

//Stop job关闭
func Stop() {
	job.Stop()
}

//Run 运行一个task
func (t *TaskRun) Run() {
	if Rds.HGet(t.User, t.TaskType).Val() != "" { //说明有相同的任务，不在执行
		return
	}
	Rds.HSet(t.User, t.TaskType, t.TaskID)

	job.AddFunc(t.Spec, t.Func)
}

//NewTaskRun 返回一个任务对象
func NewTaskRun(id, user, tasktype, spec string, function func()) *TaskRun {
	return &TaskRun{
		TaskID:   id,
		User:     user,
		TaskType: tasktype,
		Spec:     spec,
		Func:     function,
	}
}

func task() {
	job.AddFunc("", func() {})
	// job.AddJob()

}

func gettest() {}
