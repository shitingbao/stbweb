package hook

// 对mysql的日志钩子
// 资料参考 https://github.com/sohlich/elogrus
import (
	"encoding/json"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 逻辑函数结构
type hookFireFunc func(host string, entry *logrus.Entry, db *gorm.DB) error

// 钩子
// host服务地址
// client对应数据存储方
// 等级设置，设置最低记录的等级，比他高的等级都记录，包括自己
// fireFunc逻辑函数
type mysqlHook struct {
	host     string
	client   *gorm.DB
	level    logrus.Level
	fireFunc hookFireFunc
}

func (m mysqlHook) Fire(enter *logrus.Entry) error {
	return m.fireFunc(m.host, enter, m.client)
}

// 这里传入等级,只有这里定义了等级,对应等级的日志才能触发fire
func (m mysqlHook) Levels() []logrus.Level {
	return setMysqlHookLevels(m.level)
}

// NewMysqlHook,反馈一个普通hook
func NewMysqlHook(host string, client *gorm.DB, level logrus.Level) *mysqlHook {
	return newHook(host, client, level, syncFireFunc)
}

// NewAsyncMysqlHook反馈一个异步记录的hook
func NewAsyncMysqlHook(host string, client *gorm.DB, level logrus.Level) *mysqlHook {
	return newHook(host, client, level, asyncFireFunc)
}

func newHook(host string, client *gorm.DB, level logrus.Level, f hookFireFunc) *mysqlHook {
	return &mysqlHook{
		host:     host,
		client:   client,
		level:    level,
		fireFunc: f,
	}
}

func setMysqlHookLevels(level logrus.Level) []logrus.Level {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}
	return levels
}

// 异步执行
func asyncFireFunc(host string, entry *logrus.Entry, client *gorm.DB) error {
	go syncFireFunc(host, entry, client)
	return nil
}

type messgae struct {
	LogTime time.Time `gorm:"log_time"`
	Mes     string
	Data    string
	Host    string
}

// 实际逻辑操作,入库定义
func syncFireFunc(host string, entry *logrus.Entry, client *gorm.DB) error {
	da, _ := json.Marshal(entry.Data)
	m := messgae{LogTime: entry.Time, Mes: entry.Message, Data: string(da), Host: host}
	if err := client.Table("log").Select("log_time", "mes", "data", "host").Create(&m).Error; err != nil {
		log.Println("into sql:", err)
		return err
	}
	return nil
}

func HookInit(host string, client *gorm.DB, level logrus.Level) *logrus.Logger {
	return hookInit(host, client, level)
}

func hookInit(host string, client *gorm.DB, level logrus.Level) *logrus.Logger {
	lg := logrus.New()
	hk := NewMysqlHook(host, client, level)
	// hk := NewAsyncMysqlHook("", client, logrus.DebugLevel)
	lg.AddHook(hk)
	// lg.WithFields(logrus.Fields{"name": "stb"}).Info("test")
	return lg
}
