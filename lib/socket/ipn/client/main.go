package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	service         = ":1200"
	remotelyService = "124.70.156.31:1200"
)

type dataModel struct {
	User string
	Data string
}

//main监听本地请求
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError("tcp:", err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError("listen:", err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

//write和read都是阻塞形式的读写实际连接看实际应用
func handleClient(conn net.Conn) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	remotelyCon := remotelyConnect()
	go readAndWrite(&wg, conn, remotelyCon)
	go readAndWrite(&wg, remotelyCon, conn)
	wg.Wait()
}

//连接远程端口
func remotelyConnect() net.Conn {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", remotelyService)
	checkError("remotely tcp addr", err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError("remotely connect", err)
	return conn
}

//将一端的读写转发到另一端
func readAndWrite(wg *sync.WaitGroup, con, remotelyCon net.Conn) {
	for {
		request := make([]byte, 128)
		readLen, err := con.Read(request)
		if err != nil {
			log.Println("connect err:", err)
			close(remotelyCon, con)
			break
		}
		if readLen == 0 {
			log.Println("connect out")
			close(remotelyCon, con)
			break
		}
		remotelyCon.Write(request[:readLen])
	}
	wg.Done()
}
func close(remotelyCon, conn net.Conn) {
	remotelyCon.Close()
	conn.Close()
}
func checkError(reson string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, reson+": error: %s", err.Error())
		os.Exit(1)
	}
}
