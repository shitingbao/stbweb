package core

import (
	"linux_test_golang/lib/config"
	"sync"

	"github.com/shitingbao/datelogger"
)

var (
	Db            string //数据库连接
	WebConfig     config.Config
	TaskWaitGroup = new(sync.WaitGroup)
	LOG           *datelogger.DateLogger
)
