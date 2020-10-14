package loader

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stbweb/core"
	stboutserver "stbweb/lib/external_service/stb_server"
	"stbweb/lib/external_service/stbserver"
	"stbweb/lib/task"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//AutoLoader 启动项
func AutoLoader() {
	serve()

	lend := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Info("received ctrl+c,wait back job finished...")
			core.TaskWaitGroup.Wait()
			logrus.Info("all back job finished,now shutdown http server...")
			Shutdown()
			logrus.Info("success shutdown")
			lend <- true
			break
		}
	}()
	<-lend
}

func serve() {
	go func() {
		logrus.WithFields(logrus.Fields{
			"port": core.WebConfig.Port,
		}).Info("open prof")
		logrus.Info(http.ListenAndServe(fmt.Sprintf(":%s", core.WebConfig.Port), nil))
	}()
	chatHub, ctrlHub, cardHun := initChatWebsocket()
	core.Initinal(chatHub, ctrlHub, cardHun)
	// http.HandleFunc("/", httpProcess) //设置访问的路由
	clearInit()
	if core.WebConfig.ExternalServer {
		go externalServer() //开启外置服务
	}
	http.Handle("/", http.HandlerFunc(httpProcess))
}

//Shutdown 关闭所有连接
func Shutdown() {
	core.Ddb.Close()
	core.Rds.Close()
	task.Stop(core.WorkPool)
}

func externalServer() {
	lis, err := net.Listen("tcp", core.WebConfig.ExternalPort)
	if err != nil {
		logrus.Info("外置服务开启失败:", err)
		panic(err)
	}
	s := grpc.NewServer()
	stbserver.RegisterStbServerServer(s, &stboutserver.StbServe{})
	s.Serve(lis)
}
