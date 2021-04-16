package task

import (
	"encoding/json"
	"time"

	"github.com/nsqio/go-nsq"
)

// 接受逻辑，信号处理待定
// 这里改良
// 1.直接这里处理空的任务
// 2.直接在这统一反馈，不在每个操作方法中反馈
func (s *TaskElement) HandleMessage(message *nsq.Message) error {
	var (
		route   = FeedRoute
		backMes MesHandleBack
		result  = &taskMes{}
		mes     = taskFeedBack{}
	)
	mes.Sys = s.Sys
	mes.StartTime = time.Now()

	if err := json.Unmarshal(message.Body, result); err != nil {
		mes.Err = err.Error()
		mes.EndTime = time.Now()
		return sendPost(route, mes)
	}
	//空的任务或者不同版本消息，直接忽略
	if tasksList[result.TaskName] == nil || tasksList[result.TaskName].Version != result.Version {
		mes.Err = "忽略‘" + result.TaskName + "'任务"
		mes.EndTime = time.Now()
		return sendPost(route, mes)
	}

	switch result.Flag {
	case 0: // 0,开启
		route = StartRoute
		backMes = s.startTask(result.TaskName, s.Sys)
	case 1: // 1:结束
		route = EndFeedRoute
		backMes = s.stopTask(result.TaskName, s.Sys)
	case 2: // 2:更新
		route = UpdateFeedRoute
		backMes = s.updateTask(result.TaskName, s.Sys, result.Spec)
	default:
	}

	return sendPost(route, backMes)
}
