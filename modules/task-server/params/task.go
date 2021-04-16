package params

type ArgTaskRegister struct {
	Tasks   []ArgTaskBase `form:"tasks" json:"tasks"`
	Sys     string        `form:"sys" json:"sys"`
	Version float64       `form:"version" json:"version"`
}

type ArgTaskBase struct {
	TaskName string `form:"task_name" json:"task_name"`
	Spec     string `form:"spec" json:"spec"`
}

func (ArgTaskBase) TableName() string {
	return "task"
}

type ArgTask struct {
	TaskName     []string `form:"task_name" json:"task_name"`
	Sys          string   `form:"sys" json:"sys"`
	RuntimeError string   `form:"runtime_error" json:"runtime_error"`
}

type ArgTaskUpdate struct {
	ArgTask
	Spec string `form:"spec" json:"spec"`
}

// 接受
type ArgTaskMes struct {
	TaskName string  `form:"task_name" json:"task_name"`
	Spec     string  `form:"spec" json:"spec"`
	Flag     int     `form:"flag" json:"flag"` // 0,1,2:开启，结束，更新
	Version  float64 `form:"version" json:"version"`
}

type ArgTaskIssue struct {
	ArgTask
	Spec string `form:"spec" json:"spec"`
	Flag int    `form:"flag" json:"flag"`
}
