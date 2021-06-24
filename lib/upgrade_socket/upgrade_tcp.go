//这里对于长连接
//总思路，开启新进程，继承老进程的tcp服务
//老进程等待所有连接关闭后退出
//新的进程监听新的连接，老进程由于被继承不会继续监听，相当于把端口让出给新进程
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	// 设置一个重启的参数，用于区分正常开启还是升级
	argReload = flag.Bool("reload", false, "listen on fd open 3 (internal use only)")
)

// TcpCustomerListen监听入口
func TcpCustomerListen() {
	flag.Parse()
	add, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		log.Println("ResolveTCPAddr:", err)
		return
	}
	var listen net.Listener
	if *argReload {
		// 获取到cmd中的ExtraFiles内的文件信息，以它为内容启动监听
		// ExtraFiles的内容在reload方法中放入
		log.Println("*NewFile:", *argReload)
		f := os.NewFile(3, "")
		listen, err = net.FileListener(f)
		if err != nil {
			log.Println("FileListener:", err)
			return
		}

	} else {
		listen, err = net.ListenTCP("tcp", add)
		if err != nil {
			log.Println("ListenTCP:", err)
			return
		}
	}
	conCh := make(chan bool, 1)
	go func() {
		for {
			con, err := listen.Accept()
			if err != nil {
				log.Println("Accept:", err)
				return
			}
			go handle(con, conCh)
		}
	}()
	signalHandler(listen, conCh)
}

func handle(con net.Conn, conCh chan bool) {
	for {
		r := make([]byte, 256)
		n, err := con.Read(r)

		if err != nil {
			conCh <- true // 这里待定，应该适应用于多个连接
			log.Println("Read:", err)
			return
		}
		str := string(r[:n])
		log.Println(str, len(str))

		h := "333333"
		con.Write([]byte(h))
	}
}

// 信号处理
func signalHandler(listen net.Listener, conCh chan bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	log.Println("into signalHandler======")
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			log.Printf("stop")
			signal.Stop(ch)
			// listen.Close()
			log.Printf("graceful shutdown")
			return
		case syscall.SIGUSR2:
			if err := reload(listen); err != nil {
				log.Fatalf("graceful restart error: %v", err)
			}
			select {
			case <-conCh:
			}
			// listen.Close()
			log.Printf("graceful reload")
			return
		}
	}
}

// 重启方法，这里放入进程中的输入，输出和错误
// 以及最重要的listen对象（放入ExtraFiles中），以文件句柄的形式继承
// 相当于有了所有父进程的属性，然后重新执行该可执行文件
func reload(listen net.Listener) error {
	log.Println("start reload")
	f, err := listen.(*net.TCPListener).File()
	if err != nil {
		log.Println("reload", err)
		return err
	}
	cmd := exec.Command(os.Args[0], "-reload")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.ExtraFiles = append(cmd.ExtraFiles, f)
	return cmd.Start()
}

func main() {
	TcpCustomerListen()
}
