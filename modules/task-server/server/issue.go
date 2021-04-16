package server

import (
	"encoding/json"
	"net/http"
	"task-server/model"
	"task-server/params"

	"github.com/gin-gonic/gin"
	"github.com/nsqio/go-nsq"
)

var (
	tcpNsqdAddrr = "127.0.0.1:4150"
	tPro         *nsq.Producer
)

const (
	issueStart = iota
	issueEnd
	issueUpdate
)

func init() {
	config := nsq.NewConfig()
	p, err := nsq.NewProducer(tcpNsqdAddrr, config)
	if err != nil {
		panic(err)
	}
	tPro = p
}

func IssueHandle(ctx *gin.Context) {
	var err error
	arg := new(params.ArgTaskIssue)
	if err = ctx.Bind(arg); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10001,
			"msg":  err.Error(),
		})
		return
	}

	t := model.Task{}
	if err := DB.Table("task").
		Where("task_name in ? and sys = ?", arg.TaskName, arg.Sys).Find(&t).Error; err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10001,
			"msg":  err.Error(),
		})
		return
	}
	switch arg.Flag {
	case 0:
		err = IssueStart(arg.Sys, t.Version, arg.TaskName)
	case 1:
		err = IssueEnd(arg.Sys, t.Version, arg.TaskName)
	case 2:
		err = IssueUpdate(t.Version, arg)
	default:
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10001,
			"msg":  err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "ok",
	})
}

func IssueStart(topic string, version float64, name []string) error {

	for _, v := range name {
		if err := publish(topic, params.ArgTaskMes{TaskName: v, Flag: issueStart, Version: version}); err != nil {
			return err
		}
	}
	if err := DB.Table("task").
		Where("task_name in ? and sys = ?", name, topic).
		Update("is_open", 2).Debug().Error; err != nil {
		return err
	}
	return nil
}

func IssueEnd(topic string, version float64, name []string) error {
	for _, v := range name {
		if err := publish(topic, params.ArgTaskMes{TaskName: v, Flag: issueEnd, Version: version}); err != nil {
			return err
		}
	}
	return nil
}

func IssueUpdate(version float64, arg *params.ArgTaskIssue) error {
	for _, v := range arg.TaskName {
		if err := publish(arg.Sys, params.ArgTaskMes{TaskName: v, Spec: arg.Spec, Flag: issueUpdate, Version: version}); err != nil {
			return err
		}
	}
	return nil
}

func publish(topic string, arg params.ArgTaskMes) error {
	b, err := json.Marshal(arg)
	if err != nil {
		return err
	}
	return tPro.Publish(topic, b)
}
