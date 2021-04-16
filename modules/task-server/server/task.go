package server

import (
	"log"
	"net/http"
	"task-server/model"
	"task-server/params"

	"github.com/gin-gonic/gin"
)

func register(ctx *gin.Context) {
	arg := new(params.ArgTaskRegister)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	db := DB.Begin()
	tasks := []model.Task{}
	taskNames := []string{}
	for _, v := range arg.Tasks {
		t := model.Task{Sys: arg.Sys}
		t.TaskName = v.TaskName
		t.Spec = v.Spec
		t.Version = arg.Version
		tasks = append(tasks, t)
		taskNames = append(taskNames, v.TaskName)
	}

	taskReal := []model.Task{}
	if err := db.Debug().Table("task").Where("task_name in ? and sys = ?", taskNames, arg.Sys).Scan(&taskReal).Error; err != nil {
		db.Rollback()
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10001,
			"msg":  err.Error(),
		})
		return
	}

	siteTask, outTask := []model.Task{}, []model.Task{}
	for _, v := range tasks {
		isExist := true
		for _, val := range taskReal {
			if v.TaskName == val.TaskName {
				siteTask = append(siteTask, v)
				isExist = false
				break
			}
		}
		if isExist {
			outTask = append(outTask, v)
		}
	}
	for _, v := range siteTask {
		if err := db.Debug().Where("task_name = ? and sys = ?", v.TaskName, v.Sys).Updates(&v).Error; err != nil {
			db.Rollback()
			ctx.JSON(http.StatusOK, gin.H{
				"code": 10001,
				"msg":  err.Error(),
			})
			return
		}
	}
	log.Println("siteTask:", siteTask)
	log.Println("outTask:", outTask)
	if len(outTask) > 0 {
		if err := db.Debug().Create(outTask).Error; err != nil {
			db.Rollback()
			ctx.JSON(http.StatusOK, gin.H{
				"code": 10001,
				"msg":  err.Error(),
			})
			return
		}
	}

	db.Commit()
	ctx.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"msg":  "ok",
	})
}

func feedBack(ctx *gin.Context) {
	arg := new(model.TaskFeedBack)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.RuntimeError != "" {
		DB.Create(&model.TaskErr{
			TaskName:     arg.TaskName,
			Sys:          arg.Sys,
			RuntimeError: arg.RuntimeError,
		})
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "ok",
		})
		return
	}
	if err := DB.Create(arg).Error; err != nil {
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

func start(ctx *gin.Context) {
	// log.Println("into start")
	if err := taskStatusChange(1, ctx); err != nil {
		log.Println(err)
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

func stop(ctx *gin.Context) {
	if err := taskStatusChange(0, ctx); err != nil {
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

func update(ctx *gin.Context) {
	arg := new(params.ArgTaskUpdate)
	if err := ctx.Bind(arg); err != nil {
		return
	}
	if arg.RuntimeError != "" {
		errs := []model.TaskErr{}
		for _, name := range arg.TaskName {
			err := model.TaskErr{
				TaskName:     name,
				Sys:          arg.Sys,
				RuntimeError: arg.RuntimeError,
			}
			errs = append(errs, err)
		}
		DB.Create(&errs)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 10000,
			"msg":  "ok",
		})
		return
	}
	if err := DB.Table("task").
		Where("task_name in ? and sys = ?", arg.TaskName, arg.Sys).
		Update("spec", arg.Spec).Debug().Error; err != nil {
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

func taskStatusChange(isOpen int, ctx *gin.Context) error {
	// log.Println("into taskStatusChange")

	arg := new(params.ArgTask)
	if err := ctx.Bind(arg); err != nil {
		// log.Println("taskStatusChange err:", err)
		return err
	}
	if arg.RuntimeError != "" {
		errs := []model.TaskErr{}
		for _, name := range arg.TaskName {
			err := model.TaskErr{
				TaskName:     name,
				Sys:          arg.Sys,
				RuntimeError: arg.RuntimeError,
			}
			errs = append(errs, err)
		}

		return DB.Create(&errs).Error
	}
	log.Println("taskStatusChange arg:", arg)
	if err := DB.Table("task").
		Where("task_name in ? and sys = ?", arg.TaskName, arg.Sys).
		Update("is_open", isOpen).Debug().Error; err != nil {
		return err
	}
	return nil
}

// 基于被动的健康检查，清除任务
func health(ctx *gin.Context) {

}
