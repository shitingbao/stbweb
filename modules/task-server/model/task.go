package model

import "time"

type Task struct {
	TaskBase
	Sys     string
	Version float64
}

func (Task) TableName() string {
	return "task"
}

type TaskBase struct {
	TaskName string `form:"task_name" json:"task_name"`
	Spec     string `form:"spec" json:"spec"`
}
type TaskFeedBack struct {
	TaskBase
	Sys          string    `json:"sys"`
	Err          string    `form:"err" json:"err"`
	StartTime    time.Time `form:"start_time" json:"start_time"`
	EndTime      time.Time `form:"end_time" json:"end_time"`
	IsComplete   bool      `form:"is_complete" json:"is_complete"`
	RuntimeError string    `form:"runtime_error" json:"runtime_error"`
}

func (TaskFeedBack) TableName() string {
	return "task_history"
}

type TaskErr struct {
	TaskName     string `form:"task_name" json:"task_name"`
	Sys          string `json:"sys"`
	RuntimeError string `form:"runtime_error" json:"runtime_error"`
}

func (TaskErr) TableName() string {
	return "task_err"
}
