//Package task 执行一个任务
//定时任务或者普通任务
//任务中所有处理的数据都需要有日志包处理，为了防止数据丢失而设计，可以反推实现数据恢复
//任务是异步的，相互不影响，并且相同的用户提交处理相同的任务
package task

import (
	"stbweb/core"

	"github.com/pborman/uuid"

	"github.com/robfig/cron"
)

var (
	job      = cron.New()
	workPool = NewCommonPool()
)

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
	IsSave   bool   //是否保存数据包
}

func init() {
	job.Start()
}

//Stop job关闭
func Stop() {
	job.Stop()
	workPool.Release()
}

func submitPoolFunc(f func()) func() {
	return func() {
		workPool.Submit(f)
	}
}

//Run 运行一个task,内部将定时job和异步ants结合使用
func (t *Task) Run() {
	if core.Rds.HGet(t.User, t.TaskType).Val() != "" { //说明有相同的任务，不在执行
		return
	}
	core.Rds.HSet(t.User, t.TaskType, t.TaskID)

	job.AddFunc(t.Spec, submitPoolFunc(t.Func))
}

//NewTaskRun 返回一个任务对象
//user为执行用户名称，tasktype为执行类型，spec为job定时字符串，function为执行的逻辑函数
func NewTaskRun(user, tasktype, spec string, function func(), issave ...bool) *Task {
	taskID := uuid.NewUUID().String()
	sa := true
	if len(issave) > 0 {
		sa = issave[0]
	}
	return &Task{
		TaskID:   taskID,
		User:     user,
		TaskType: tasktype,
		Spec:     spec,
		Func:     function,
		IsSave:   sa,
	}
}

// func task() {
// 	job.AddFunc("", func() {})
// 	// job.AddJob()

// }
