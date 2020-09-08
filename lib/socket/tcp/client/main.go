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

//write和read都是阻塞形式的读写实际连接看实际应用
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError("tcp addr", err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError("connect", err)
	for {
		result := make([]byte, 256)
		readLen, err := conn.Read(result)
		checkError("conn read", err)
		res := dataModel{}
		//这里的json解析得去掉后面的空格，不然会有invalid character '\x00' after top-level value，这个错误，可能是结尾有什么特殊字符导致的
		if err := json.Unmarshal(result[:readLen], &res); err != nil {
			checkError("unmarshal:", err)
			break
		}
		log.Println("get mes:", string(result))
		da, _ := json.Marshal(dataModel{User: "shitinbao", Data: "123"})
		wlen, err := conn.Write(da)
		checkError("conn write", err)
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
