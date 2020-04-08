package common

import (
	"stbweb/core"
	"stbweb/lib/rediser"
	"stbweb/lib/task"
)

func init() {
	ts := task.NewTask("sys", "clearMember", "0 0/5 * * * ?", clearFun)
	ts.Run()
}

var clearFun = func() {
	rediser.ClearMember(core.Rds)
}
