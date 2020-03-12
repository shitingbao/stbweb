package core

import (
	"linux_test_golang/lib/config"
	"sync"

	"github.com/shitingbao/datelogger"
)

var (
	//Db 数据库连接
	Db string

	//WebConfig 数据库连接
	WebConfig config.Config

	//TaskWaitGroup 任务
	TaskWaitGroup = new(sync.WaitGroup)

	//LOG 日志
	LOG *datelogger.DateLogger
)
