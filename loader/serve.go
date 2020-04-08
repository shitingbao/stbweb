package loader

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"stbweb/core"
	"stbweb/lib/task"

	"github.com/Sirupsen/logrus"
)

//AutoLoader 启动项
func AutoLoader() {
	serve()

	lend := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			core.LOG.Info("received ctrl+c,wait back job finished...")
			core.TaskWaitGroup.Wait()
			core.LOG.Info("all back job finished,now shutdown http server...")
			Shutdown()
			core.LOG.Info("success shutdown")
			lend <- true
			break
		}
	}()
	<-lend
}

func serve() {
	go func() {
		core.LOG.WithFields(logrus.Fields{
			"port": core.WebConfig.Port,
		}).Info("open prof")
		core.LOG.Info(http.ListenAndServe(fmt.Sprintf(":%s", core.WebConfig.Port), nil))
	}()
	chatHub, ctrlHub := initChatWebsocket()
	core.Initinal(chatHub, ctrlHub)
	// http.HandleFunc("/", httpProcess) //设置访问的路由
	clearInit()
	http.Handle("/", http.HandlerFunc(httpProcess))
}

//Shutdown 关闭所有连接
func Shutdown() {
	core.Ddb.Close()
	core.Rds.Close()
	task.Stop()
}
