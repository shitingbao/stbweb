package common

import (
	"stbweb/core"
	"stbweb/lib/rediser"
	"stbweb/lib/task"
)

func init() {
	ts := task.NewTask("sys", "clearMember", "*/5 * * * ?", clearFun) //五分钟执行一次
	ts.Run()
}

var clearFun = func() {
	rediser.ClearMember(core.Rds)
}
