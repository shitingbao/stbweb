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
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	// set maxium request length to 128B to prevent flood attack
	defer conn.Close() // close connection before exit
	for {
		conn.Write([]byte("shitngbao"))
		log.Println("start read")
		request := make([]byte, 128)
		readLen, err := conn.Read(request)
		log.Println("read len:", readLen)
		if err != nil {
			fmt.Println("connect:", err)
			break
		}
		if readLen == 0 {
			break // connection already closed by client
		}
		log.Println("client get mes:", string(request))

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
