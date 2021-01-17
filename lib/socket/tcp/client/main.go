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
	Data int
}

//write和read都是阻塞形式的读写实际连接看实际应用
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError("tcp addr:", err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError("connect:", err)
	_, err = conn.Write([]byte("ep-ctl:15164350934"))
	checkError("Write:", err)
	res := make([]byte, 256)
	ln, err := conn.Read(res)
	checkError("read", err)
	if string(res[:ln]) == "success" {
		log.Println("success")
	} else {
		return
	}
	for {
		i := 0
		da, _ := json.Marshal(dataModel{User: "shitinbao", Data: i})
		wlen, err := conn.Write(da)
		if err != nil {
			log.Println("write:", err)
			return
		}
		checkError("Write:", err)
		log.Println("发送长度:", wlen)
		// result, err = ioutil.ReadAll(conn)
		time.Sleep(time.Second * 2)
	}
}
func checkError(reson string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, reson+": %s", err.Error())
		os.Exit(1)
	}
}
