package core

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"stbweb/lib/chatroom"
	"stbweb/lib/config"
	"stbweb/lib/ddb"
	"stbweb/lib/mongodb"
	"stbweb/lib/rediser"
	"stbweb/lib/task"
	"stbweb/lib/ws"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/panjf2000/ants"
	"github.com/sirupsen/logrus"
)

var (
	//DefaultFilePath 输出默认路径
	DefaultFilePath = "./assets"

	//Ddb 数据库连接
	Ddb *sql.DB

	//WebConfig config配置
	WebConfig *config.Config

	//TaskWaitGroup 任务
	TaskWaitGroup = new(sync.WaitGroup)

	//LOG 日志
	// LOG *datelogger.DateLogger

	//ChatHub 公共聊天频道的hub对象
	ChatHub *ws.Hub

	//CtrlHub 发送控制消息的hub对象
	CtrlHub *ws.Hub
	//CardHun 牌hub对象
	CardHun *ws.Hub
	//RoomChatHub 新聊天室对象
	RoomChatHub *RoomChatHubSet
	//Rds redis连接d
	Rds *redis.Client
	//WorkPool 全局工作池
	WorkPool *ants.Pool

	//Mdb mongodb连接对象
	Mdb *mongodb.Mongodb

	//RoomLocks 房间号对应的锁结构，key为roomid
	RoomLocks = make(map[string]*CustomizeLock)

	//RoomSets 房间集合对象，Romid对应一个room,保存所有的房间唯一号和房间对象的对应关系
	RoomSets = make(map[string]*chatroom.ChatRoom)
)

func init() {
	checkLog()
	WorkPool = task.JobInit()
}
func checkConfig() {
	WebConfig = config.ReadConfig("./config.json") //配置准备
}

//初始化日志文件，如果已经初始化则跳过,并获取配置参数
//重定向日志输出的文件
func checkLog() {
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
	// LOG = &datelogger.DateLogger{Path: filepath.Join(workDir, "log"), Level: lvl}
	// }
}

//Initinal 函数初始化日志及数据库链接，以及以后的消息频道
func Initinal(chatHub, ctrlHub, cardHun *ws.Hub, roomChatHub *RoomChatHubSet) {
	ChatHub = chatHub
	CtrlHub = ctrlHub
	CardHun = cardHun
	RoomChatHub = roomChatHub
	pathExists()
	if err := openx(WebConfig.Driver, WebConfig.ConnectString); err != nil {
		logrus.WithFields(logrus.Fields{"Driver": WebConfig.Driver, "ConnectString": WebConfig.ConnectString}).Panic("database")
		// LOG.Printf("open database error drive %s ,connection string:%s\n", WebConfig.Driver, WebConfig.ConnectString)
	}
	openRdis(WebConfig.RedisAdree+":"+WebConfig.RedisPort, WebConfig.RedisPwd, WebConfig.Redislevel)
	// restoreConnect()
	openMongodb(WebConfig.MongoDriver, WebConfig.MongoDatabase)
	return
}

func openMongodb(driver, database string) {
	mongo, err := mongodb.OpenMongoDb(driver, database)
	if err != nil {
		logrus.WithFields(logrus.Fields{"ConnectString": driver, "database": database, "err:": err}).Panic("mongodb err")
	}
	logrus.WithFields(logrus.Fields{"Connect": driver + "/" + database}).Info("mongodb")
	Mdb = mongo
}

//Openx 打开一个数据库连接，返回一个包装过的DB对象，其能返回DriverName
func openx(driverName, dataSourceName string) error {
	defer func() {
		if err := recover(); err != nil {
			logrus.Info("open db have err:", err)
			logrus.Info(driverName, ":", dataSourceName, "--db连接5S后重试。。。。。。")
		}
	}()
	d, err := ddb.Open(driverName, dataSourceName)
	if err != nil {
		logrus.Info("open db have err:", err)
		logrus.Info(driverName, ":", dataSourceName, "--db连接5S后重试。。。。。。")
		return err
	}
	Ddb = d
	logrus.Info("Driver:", WebConfig.Driver, "--ConnectString:", WebConfig.ConnectString, "  connect success!")
	return nil
}

func openRdis(addr, pwd string, dbevel int) {
	defer func() {
		if err := recover(); err != nil {
			logrus.WithFields(logrus.Fields{"ConnectString": addr, "level": dbevel}).Panic("redis")
		}
	}()
	Rds = rediser.Open(addr, pwd, dbevel)
}

//pathExists 判断是否存在默认路径，不存在则生成
func pathExists() {
	_, err := os.Stat(DefaultFilePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{"msg": err.Error()}).Error("DefaultFilePath")
	}
	if os.IsNotExist(err) {
		if err := os.MkdirAll(DefaultFilePath, os.ModePerm); err != nil {
			logrus.WithFields(logrus.Fields{"msg": err.Error()}).Error("CreateDefaultFilePath")
		} else {
			logrus.Info("生成默认路径")
		}
	}
}

//restoreConnect db连接重启,5s后重启,暂时保留
func restoreConnect() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		logrus.Info("DB重连start")
		for range ticker.C {
			if Ddb == nil {
				logrus.Info("DB正在重新连接。。。。。。")
				openx(WebConfig.Driver, WebConfig.ConnectString)
			}
			if Rds == nil {
				logrus.Info("redis正在重新连接。。。。。。")
				openRdis(WebConfig.RedisAdree+":"+WebConfig.RedisPort, WebConfig.RedisPwd, WebConfig.Redislevel)
			}
		}
	}()
	// defer ticker.Stop()
}
