package core

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"stbweb/lib/config"
	"stbweb/lib/ddb"
	"stbweb/lib/rediser"
	"stbweb/lib/ws"
	"sync"

	"github.com/Sirupsen/logrus"
	sysRedis "github.com/go-redis/redis"
	"github.com/shitingbao/datelogger"
)

var (
	//DefaultFilePath 输出默认路径
	DefaultFilePath = "./assets"

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

	//Rds redis连接d
	Rds *sysRedis.Client
)

func init() {
	checkLog()
}
func checkConfig() {
	WebConfig = config.ReadConfig("./config.json") //配置准备
}

//初始化日志文件，如果已经初始化则跳过,并获取配置参数
//重定向日志输出的文件
func checkLog() {
	if LOG == nil {
		checkConfig()
		str, err := os.Executable()
		if err != nil {
			logrus.Panic(err)
		}
		workDir := filepath.Dir(str)
		flog, err := os.OpenFile(filepath.Join(workDir, "log.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		//将所有的panic信息输出到err文件中，而不是控制台，因为控制台有行数限制
		//https://stackoverflow.com/questions/34772012/capturing-panic-in-golang
		//直接使用logrus输出，将输出在当前的log.txt文件中
		//使用LOG对象输出，将输出在log目录下的每日日志当中，详细看对应文件句柄对照
		//panic将输出在err文件中，期间调用底层panic重定向方法
		ferr, err := os.OpenFile(filepath.Join(workDir, "err.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Errorf("error opening file: %v", err))
		}
		redirectStderr(ferr)
		log.SetFlags(log.LstdFlags | log.Llongfile)
		logrus.SetOutput(io.MultiWriter(os.Stdout, flog))
		lvl, err := logrus.ParseLevel(WebConfig.LogLevel)
		if err != nil {
			panic(err)
		}
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
		logrus.SetLevel(lvl)
		logrus.WithFields(logrus.Fields{"set-level": lvl.String()}).Info("initlog")

		LOG = &datelogger.DateLogger{Path: filepath.Join(workDir, "log"), Level: lvl}
	}
}

//Initinal 函数初始化日志及数据库链接，以及以后的消息频道
func Initinal(chatHub, ctrlHub *ws.Hub) {
	ChatHub = chatHub
	CtrlHub = ctrlHub
	pathExists()
	if err := openx(WebConfig.Driver, WebConfig.ConnectString); err != nil {
		LOG.WithFields(logrus.Fields{"Driver": WebConfig.Driver, "ConnectString": WebConfig.ConnectString}).Panic("database")
		// LOG.Printf("open database error drive %s ,connection string:%s\n", WebConfig.Driver, WebConfig.ConnectString)
	}
	openRdis(WebConfig.RedisAdree, WebConfig.RedisPwd, WebConfig.Redislevel)
	return
}

//Openx 打开一个数据库连接，返回一个包装过的DB对象，其能返回DriverName
func openx(driverName, dataSourceName string) error {
	d, err := ddb.Open(driverName, dataSourceName)
	if err != nil {
		return err
	}
	Ddb = d
	return nil
}
func openRdis(addr, pwd string, dbevel int) {
	Rds = rediser.Open(addr, pwd, dbevel)
	rediser.SetUser(Rds)
	rediser.GetUser(Rds)

}

//pathExists 判断是否存在默认路径，不存在则生成
func pathExists() {
	_, err := os.Stat(DefaultFilePath)
	if err != nil {
		LOG.WithFields(logrus.Fields{"msg": err.Error()}).Error("DefaultFilePath")
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(DefaultFilePath, os.ModePerm); err != nil {
			LOG.WithFields(logrus.Fields{"msg": err.Error()}).Error("DefaultFilePath")
		}
	}
}
