package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	service = ":1200"
)

type dataModel struct {
	User string
	Data string
}

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
	log.Println("adress:", conn.LocalAddr().String())
	defer conn.Close() // close connection before exit
	for {

		da, _ := json.Marshal(dataModel{User: "shitinbao", Data: "123"})
		conn.Write(da)
		request := make([]byte, 128)
		readLen, err := conn.Read(request)
		if readLen == 0 {
			log.Println("connect out")
			break // connection already closed by client
		}
		res := dataModel{}
		if err := json.Unmarshal(request[:readLen], &res); err != nil {
			log.Println("err:", err)
			break
		}
		log.Println("get client:", string(request), "-接受长度:", readLen)
		if err != nil {
			fmt.Println("connect:", err)
			break
		}

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
