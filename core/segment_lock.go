//Package core 使用该锁注意，依赖于redis
//1.一个get，对应一个Freed，获取和释放需要对应，不然通道内的信号量达不到一致性可能会
//2.在使用时，应该给他设置以user为key的setnx
package core

//CustomizeLock cap总长度，锁实际内容，key标识所属（user）
//flag标识过程的锁，用户限制连接数
type CustomizeLock struct {
	cap   int
	locks map[string]bool
	flag  chan bool
}

//NewLock 新建一个锁，放入需要使用的长度,以及对应roomid的标识
//长度应当不能少于2，对应房间人数
//同时在flag队列里面放满标识待用
//注意这里不能关闭标识通道，因为需要复用，重新加入标识（比如使用完毕或者超时，重新放入，等待下一次使用）
//这里生成后直接加入对应全局房间锁列表中
func NewLock(cap int, roomID string) *CustomizeLock {
	flag := make(chan bool, cap)
	for i := 0; i < cap; i++ {
		flag <- true
	}
	ck := &CustomizeLock{
		cap:   cap,
		locks: make(map[string]bool),
		flag:  flag,
	}
	RoomLocks[roomID] = ck
	// close(flag)
	return ck
}

//GetLock 从标识队列中获取一个锁,并加入使用对象
//关闭通道说明房间移除，临界情况为，过程中加入丽连接
func (c *CustomizeLock) GetLock(user string) bool {
	select {
	case _, ok := <-c.flag:
		if !ok {
			return false
		}
		c.locks[user] = true
		return true
	default:
		return false
	}
}

// FreedLock 释放锁
func (c *CustomizeLock) FreedLock(user string) {
	delete(c.locks, user)
	c.flag <- true //记得放回标识
}

//Clear 清理锁
func (c *CustomizeLock) Clear(roomID string) {
	close(c.flag)
	delete(RoomLocks, roomID)
}
