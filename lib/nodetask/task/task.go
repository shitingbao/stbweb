package task

// 适应与
import (
	"encoding/json"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/robfig/cron/v3"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.0.1:4150"
	jobsList     = make(map[string]*TaskHandler)
)

type TaskBase struct {
	TaskName string
	Spec     string
}

// 接受
type TaskMes struct {
	TaskBase
	Flag int // 0,1,2:开启，结束，更新
}

// 继承，并实现task接口
type TaskHandler struct {
	TaskBase
	EntryID cron.EntryID
	Func    func() error //需要执行的逻辑
}

type TaskElement struct {
	Topic   string
	Channel string
	Sys     string
	Handers []*TaskHandler
	Cron    *cron.Cron
}
type taskFeedBack struct {
	TaskBase
	Err        string
	EndTime    time.Time
	IsComplete bool
}

// GetSymbolize反馈对应的主题和通道，为了能直接在开启监听中直接获取
func (t *TaskElement) GetSymbolize() (string, string) {
	return t.Topic, t.Channel
}

func NewTaskHandle(taskName, spec string, fc func() error) *TaskHandler {
	t := new(TaskHandler)
	t.TaskName = taskName
	t.Spec = spec
	t.Func = fc
	return t
}

func NewTaskElement(topic, channel, sys string, c *cron.Cron, task ...*TaskHandler) *TaskElement {
	return &TaskElement{
		Topic:   topic,
		Channel: channel,
		Sys:     sys,
		Handers: task,
		Cron:    c,
	}
}

// 包装逻辑执行，加上错误处理和反馈
func (t *TaskHandler) jobFun(fc func() error) func() {
	return func() {
		mes := taskFeedBack{}
		mes.TaskBase = t.TaskBase
		mes.IsComplete = true
		if err := fc(); err != nil {
			mes.Err = err.Error()
			mes.IsComplete = false
		}
		//发送反馈
		feedBack(mes)
	}
}

// 接受逻辑，信号处理待定
func (s *TaskElement) HandleMessage(message *nsq.Message) error {
	result := &TaskMes{}
	if err := json.Unmarshal(message.Body, result); err != nil {
		return err
	}
	switch result.Flag {
	case 0: // 0,开启
		s.StartJob(result.TaskName)
	case 1: // 1:结束
		s.StopJob(result.TaskName)
	case 2: // 2:更新
		s.UpdateJob(result.TaskName, result.Spec)
	default:
	}
	return nil
}

// StartJob应该先检查是否存在该任务
func (s *TaskElement) StartJob(jobName string) error {
	job := jobsList[jobName]
	if job.EntryID != 0 {
		return nil
	}
	id, err := s.Cron.AddFunc(job.Spec, job.jobFun(job.Func))
	if err != nil {
		return err
	}
	job.EntryID = id
	return nil
}

func (s *TaskElement) StopJob(jobName string) {
	job := jobsList[jobName]
	job.EntryID = 0
	s.Cron.Remove(job.EntryID)
}

func (s *TaskElement) UpdateJob(jobName, spec string) error {
	job := jobsList[jobName]
	s.Cron.Remove(job.EntryID)
	id, err := s.Cron.AddFunc(spec, job.jobFun(job.Func))
	if err != nil {
		return err
	}
	jobsList[jobName].EntryID = id
	jobsList[jobName].Spec = spec
	return nil
}

// 设置并发送task
func setJobs(task *TaskElement) error {
	body := []TaskBase{}
	for _, v := range task.Handers {
		jobsList[v.TaskName] = v
		body = append(body, v.TaskBase)
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
	if err := setJobs(task); err != nil {
		return err
	}
	consumer.AddHandler(task)
	if err := consumer.ConnectToNSQD(tcpNsqdAddrr); err != nil {
		return err
	}
	return nil
}
