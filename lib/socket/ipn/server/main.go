package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	service = ":1200"
)

var (
	connLists  map[net.Conn]bool
	connectNum = 0
)

//socket的数据类型
type dataModel struct {
	User string
	Data string
}

//过程1，接受两个连接
//过程2，将1个数据读取后放入另一个连接中去
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
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
	connLists[conn] = true
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	// log.Println("adress:", conn.LocalAddr().String())
	defer close(conn)
	for {
		request := make([]byte, 128)
		readLen, err := conn.Read(request)
		if err != nil {
			log.Println("read err:", err)
			delete(connLists, conn)
			break
		}
		if readLen == 0 {
			log.Println("connect 断开连接")
			delete(connLists, conn)
			break
		}
		// res := dataModel{}
		// if err := json.Unmarshal(request[:readLen], &res); err != nil {
		// 	log.Println("err:", err)
		// 	delete(connLists, conn)
		// 	break
		// }
		// log.Println("get client:", string(request), "-接受长度:", readLen)
		// if err != nil {
		// 	log.Println("connect:", err)
		// 	delete(connLists, conn)
		// 	break
		// }
		broadcast(conn, request)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//广播所有除自己外的所有连接
func broadcast(tcon net.Conn, da []byte) {
	for con := range connLists {
		if tcon == con {
			continue
		}
		con.Write(da)
	}
}
func close(conn net.Conn) {
	conn.Close()
	delete(connLists, conn)
	log.Println("connect 断开连接")
}
