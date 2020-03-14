package core

import (
	"database/sql"
	"fmt"
	"io"
	syslog "log"
	"os"
	"path/filepath"
	"stbweb/lib/config"
	"stbweb/lib/ws"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/shitingbao/datelogger"
)

var (
	//Ddb 数据库连接
	Ddb *sql.DB

	//WebConfig 数据库连接
	WebConfig *config.Config

	//TaskWaitGroup 任务
	TaskWaitGroup = new(sync.WaitGroup)

	//LOG 日志
	LOG *datelogger.DateLogger
	//ChatHub 公共聊天频道的hub对象
	ChatHub *ws.Hub

	//CtrlHub 发送控制消息的hub对象
	CtrlHub *ws.Hub
)

func checkConfig() {
	WebConfig = config.ReadConfig("./config.json") //配置准备
}

//初始化日志文件，如果已经初始化则跳过,并获取配置参数
func checkLog() {
	if LOG == nil {
		checkConfig()
		str, err := os.Executable()
		if err != nil {
			log.Panic(err)
		}
		workDir := filepath.Dir(str)
		flog, err := os.OpenFile(filepath.Join(workDir, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		//将所有的panic信息输出到err文件中，而不是控制台，因为控制台有行数限制
		//https://stackoverflow.com/questions/34772012/capturing-panic-in-golang
		ferr, err := os.OpenFile(filepath.Join(workDir, "err.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		redirectStderr(ferr)
		syslog.SetFlags(syslog.LstdFlags | syslog.Llongfile)
		log.SetOutput(io.MultiWriter(os.Stdout, flog))
		lvl, err := log.ParseLevel(WebConfig.LogLevel)
		if err != nil {
			panic(err)
		}
		log.SetLevel(lvl)
		log.WithFields(log.Fields{"set-level": lvl.String()}).Info("initlog")
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: "20060102T150405",
		})
		LOG = &datelogger.DateLogger{Path: filepath.Join(workDir, "log"), Level: lvl}
	}
}

//Initinal 函数初始化日志及数据库链接，以及以后的消息频道
func Initinal(chatHub, ctrlHub *ws.Hub) {
	ChatHub = chatHub
	CtrlHub = ctrlHub
	checkLog()
	if err := openx(WebConfig.Driver, WebConfig.ConnectString); err != nil {
		LOG.Printf("open database drive %s ,connection string:%s\n", WebConfig.Driver, WebConfig.ConnectString)
	}
	return
}

//Openx 打开一个数据库连接，返回一个包装过的DB对象，其能返回DriverName
func openx(driverName, dataSourceName string) error {
	d, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	Ddb = d
	return nil
}
