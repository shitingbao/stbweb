package task

// 适应与
import (
	"errors"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.0.1:4150"
	tasksList    = make(map[string]*TaskHandler)
)

// runtimeErr 应该都有一个err反馈
type runtimeErr struct {
	RuntimeError string `json:"runtime_error"`
}
type argTask struct {
	runtimeErr
	TaskName []string `json:"task_name"`
	Sys      string   `json:"sys"`
}
type argTaskUpdate struct {
	argTask
	runtimeErr
	Spec string `json:"spec"`
}
type argRegister struct {
	Tasks   []taskBase `json:"tasks"`
	Sys     string     `json:"sys"`
	Version float64    `json:"version"`
}

type taskBase struct {
	TaskName string `json:"task_name"`
	Spec     string `json:"spec"`
}

// 接受
type taskMes struct {
	taskBase
	Flag    int     `json:"flag"` // 0,1,2:开启，结束，更新
	Version float64 `json:"version"`
}

// 继承，并实现task接口
type TaskHandler struct {
	taskBase
	EntryID cron.EntryID
	Func    func() error //需要执行的逻辑
	Version float64      `json:"version"`
}

type TaskElement struct {
	Topic   string
	Channel string
	Sys     string
	Version float64
	Handers []*TaskHandler
	Cron    *cron.Cron
}

type taskFeedBack struct {
	taskBase
	Sys        string    `json:"sys"`
	Err        string    `json:"err"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	IsComplete bool      `json:"is_complete"`
}

func NewTaskHandle(taskName, spec string, fc func() error) *TaskHandler {
	t := new(TaskHandler)
	t.TaskName = taskName
	t.Spec = spec
	t.Func = fc
	return t
}

func NewTaskElement(topic, channel, sys string, version float64, task ...*TaskHandler) *TaskElement {
	c := cron.New()
	c.Start()
	return &TaskElement{
		Topic:   topic,
		Channel: channel,
		Sys:     sys,
		Version: version,
		Cron:    c,
		Handers: task,
	}
}

// 包装逻辑执行，加上错误处理和反馈
func (t *TaskHandler) taskFun(sys string, fc func() error) func() {
	return func() {
		mes := taskFeedBack{}
		mes.taskBase = t.taskBase
		mes.IsComplete = true
		mes.Sys = sys
		mes.StartTime = time.Now()

		if err := fc(); err != nil {
			mes.Err = err.Error()
			mes.IsComplete = false
		}
		mes.EndTime = time.Now()
		//发送反馈
		feedBack(mes)
	}
}

// startTask应该先检查是否存在该任务
func (s *TaskElement) startTask(taskName, sys string) argTask {
	arg := argTask{
		TaskName: []string{taskName},
		Sys:      sys,
	}
	// defer start(arg)
	task := tasksList[taskName]

	if task == nil || task.EntryID != 0 { //说明已经开启了
		err := errors.New(taskName + " task start have err")
		arg.RuntimeError = err.Error()
		return arg
	}

	id, err := s.Cron.AddFunc(task.Spec, task.taskFun(s.Sys, task.Func))
	if err != nil {
		arg.RuntimeError = err.Error()
		return arg
	}
	task.EntryID = id

	return arg
}

func (s *TaskElement) stopTask(taskName, sys string) argTask {
	arg := argTask{
		TaskName: []string{taskName},
		Sys:      sys,
	}
	// defer end(arg)
	task := tasksList[taskName]
	if task == nil || task.EntryID == 0 {
		arg.RuntimeError = taskName + " task end have err"
		return arg
	}
	s.Cron.Remove(task.EntryID)
	// delete(tasksList, taskName)//不能直接删除，这样就没这个对象了，id清空即可
	tasksList[taskName].EntryID = 0
	return arg
}

func (s *TaskElement) updateTask(taskName, sys, spec string) argTaskUpdate {
	arg := argTaskUpdate{}
	arg.TaskName = []string{taskName}
	arg.Sys = sys
	arg.Spec = spec
	// defer update(arg)

	task := tasksList[taskName]
	if task == nil {
		err := errors.New(taskName + " this task is nil")
		arg.RuntimeError = err.Error()
		return arg
	}
	s.Cron.Remove(task.EntryID)
	id, err := s.Cron.AddFunc(spec, task.taskFun(s.Sys, task.Func))
	if err != nil {
		arg.RuntimeError = err.Error()
		return arg
	} //注意，先增加在删除

	tasksList[taskName].EntryID = id
	tasksList[taskName].Spec = spec

	return arg
}

func (s *TaskElement) Close() {
	arg := argTask{
		Sys: s.Sys,
	}
	defer func() {
		end(arg)
		s.Cron.Stop()
	}()
	taskNames := []string{}
	for key, task := range tasksList {
		taskNames = append(taskNames, key)
		s.Cron.Remove(task.EntryID)
	}
	arg.TaskName = taskNames
	tasksList = make(map[string]*TaskHandler)
}

// 设置并发送task
func setTasks(task *TaskElement) error {
	body := argRegister{Sys: task.Sys, Version: task.Version}
	for _, v := range task.Handers {
		v.Version = task.Version
		tasksList[v.TaskName] = v
		body.Tasks = append(body.Tasks, v.taskBase)
	}
	return register(body)
}

// 开启观察
func Watch(task *TaskElement) error {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(task.Topic, task.Channel, config)
	if err != nil {
		return err
	}

	consumer.AddHandler(task)
	if err := consumer.ConnectToNSQD(tcpNsqdAddrr); err != nil {
		return err
	}
	return setTasks(task)
}
