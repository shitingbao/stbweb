package loader

import (
	"stbweb/core"
	"stbweb/lib/rediser"
	"stbweb/lib/task"
)

func clearInit() {
	ts := task.NewTask("sys", "clearMember", "*/5 * * * ?", clearFun) //五分钟执行一次
	ts.Run()
}

var clearFun = func() error {
	rediser.ClearMember(core.Rds)
	return nil
}
