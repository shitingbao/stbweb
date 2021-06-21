package upgrade

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
	var listen net.Listener
	add, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		log.Println("ResolveTCPAddr:", err)
		return
	}
	if *argReload {
		// 获取到cmd中的ExtraFiles内的文件信息，以它为内容启动监听
		// ExtraFiles的内容在reload方法中放入
		f := os.NewFile(3, "")
		listen, err = net.FileListener(f)
		log.Println("FileListener:", err)

	} else {
		listen, err = net.ListenTCP("tcp", add)
		log.Println("ListenTCP:", err)
	}
	go func() {
		for {
			con, err := listen.Accept()
			if err != nil {
				log.Println("Accept:", err)
				return
			}
			go handle(con)
		}
	}()
	signalHandler(listen)
}

func handle(con net.Conn) {
	for {
		r := make([]byte, 256)
		if _, err := con.Read(r); err != nil {
			log.Println("Read:", err)
			return
		}
		log.Println(string(r))
		con.Write([]byte("hello world"))
	}
}

// 信号处理
func signalHandler(listen net.Listener) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// stop
			log.Printf("stop")
			signal.Stop(ch)
			listen.Close()
			log.Printf("graceful shutdown")
			return
		case syscall.SIGUSR2:
			if err := reload(listen); err != nil {
				log.Fatalf("graceful restart error: %v", err)
			}
			listen.Close()
			log.Printf("graceful reload")
			return
		}
	}
}

// 重启方法，这里放入进程中的输入，输出和错误
// 以及最重要的listen对象（放入ExtraFiles中），以文件句柄的形式继承
// 相当于有了所有父进程的属性，然后重新执行该可执行文件
func reload(listen net.Listener) error {
	f, err := listen.(*net.TCPListener).File()
	if err != nil {
		return err
	}
	cmd := exec.Command(os.Args[0], "-reload")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.ExtraFiles = append(cmd.ExtraFiles, f)
	return cmd.Start()
}
