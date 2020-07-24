package core

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants"
)

const (
	//测试并行数量
	runTime = 100

	//BenchAntsSize 同上
	BenchAntsSize = 200000
	//DefaultExpiredTime 测试使用默认过期时间
	DefaultExpiredTime = 10 * time.Second
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demosFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

//antsCommon 使用普通的pool
//例子一
func antsCommon() {
	p, _ := ants.NewPool(BenchAntsSize) //新建一个pool对象，其他同上
	defer p.Release()
	for j := 0; j < runTime; j++ {
		_ = p.Submit(func() {
			log.Println(":hello")
			time.Sleep(time.Millisecond * 10)
		})
	}

}

//antsMarkFuncPut 使用特定的带有函数内容的pool
//例子二
func antsMarkFuncPut() {
	// var wg sync.WaitGroup

	p, _ := ants.NewPoolWithFunc(10, func(i interface{}) { //新建一个带有同类方法的pool对象
		myFunc(i)
		// wg.Done()
	})
	defer p.Release()
	// Submit tasks one by one.
	for i := 0; i < 1000; i++ {
		// wg.Add(1)
		_ = p.Invoke(int32(i)) //这个就是发送，相当于上述普通的pool的submit，唯一不同的是参数，因为这个发送的同类型的方法，加入逻辑代码时注意
	}
	// wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
	// if sum != 499500 {
	// 	panic("the final result is wrong!!!")
	// }
}

//antsDefaultCommon 使用默认普通pool
//其实就是使用了普通的pool，为了方便直接使用，在内部已经new了一个普通的pool，
//相当于下面那个新建的过程给你写好了，容量大小和过期时间都用默认的，详细信息可以看源码，里面剥一层就可以看到
//例子三
func antsDefaultCommon() {
	// var wg sync.WaitGroup //这里使用等待是为了看出结果，阻塞主线程，防止直接停止，如果在web项目中，就不需要

	defer ants.Release() //退出工作，相当于使用后关闭
	log.Println("start ants work")
	for i := 0; i < runTime; i++ {
		// wg.Add(1)
		ants.Submit(func() { //提交函数，将逻辑函数提交至work中执行，这里写入自己的逻辑
			log.Println(i, ":hello")
			// time.Sleep(time.Millisecond * 10)
			// wg.Done()
		})
	}
	// wg.Wait()
	log.Println("stop ants work")
}

//注意
//线程池执行有两种，一种执行普通逻辑方法pool，可接受所有方法,另一种执行形同类型的方法（就是每次接收的内容方法都一样）
//使用前需要先建立一个对应的pool对象，参数是容量大小和过期时间等， 如果使用普通方法，有默认方法可以使用
//使用后，需要调用Release来结束使用，意思就是close那类的意思
//总的来说，使用过程就是，新建对象，加入逻辑函数代码，结束使用关闭对象。

//在实际使用过程中，应该对不同逻辑新建不同的pool，这样处理后关闭，相互不会影响
