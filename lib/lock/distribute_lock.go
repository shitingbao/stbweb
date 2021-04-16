package lock

// 两个生产者，一个消费者的基本结构
// 退出的逻辑，续命和退出使用同一个通道，接受到后判断该信号是哪种类型来进行操作
// 这样避免select同时接受信号时的数据丢失
import (
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	lockNx   = "log_lock"
	baseTime = time.Second * 50
	outTime  = time.Second * 40
)

// GetDistributeLock简易带守护进程的分布式锁,
// 用户，执行逻辑函数，成功执行反馈true
// user标识谁在使用这把锁，fc为该锁执行的逻辑，name对应user
// 因为1的信号是close之前的，所以执行close之前肯定有 的1（删除对应key标记）操作，所以执行0的时候判断key即可，不用担心通道close
func GetDistributeLock(user string, Rds *redis.Client, fc func(user string)) bool {
	//锁的标识得是有一个身份辨认,锁获取成功才继续
	ok, err := Rds.SetNX(lockNx, user, baseTime).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{"setnx": err.Error()}).Error("redis")
		return false
	}
	if !ok {
		return false
	}
	out := make(chan int) //0代表续命，1代表死亡
	go func() {           //续命携程
		for {
			select {
			case <-time.After(outTime):
				if Rds.Get(lockNx).Val() != user { //只有自己的锁，才能给他续命//保险操作
					//这时候该用户的操作已经退出来，所以应该直接退出
					//@double check 1，1，2对应
					return
				}
				out <- 0
			}
		}
	}()
	go func() { //清理锁的携程
		for {
			select {
			case sign := <-out:
				switch sign {
				case 0:
					if Rds.Get(lockNx).Val() != user { //只有自己的锁，才能给他续命//保险操作
						//@double check 2，1，2对应
						return
					}
					Rds.Expire(lockNx, time.Minute)
				case 1:
					Rds.Del(lockNx) //成功后释放锁
					return
				}
			}
		}
	}()
	defer func() { //fc错误收集
		if err := recover(); err != nil {
			go func() { out <- 1; close(out) }() //说明执行fc异常了，释放锁
			logrus.WithFields(logrus.Fields{"Execution function": err}).Error("DistributeLock")
		}
	}()
	fc(user)
	go func() { out <- 1; close(out) }() //直接退出守护进程,这里要关闭out，防止select中两个同时到期
	return true
}
