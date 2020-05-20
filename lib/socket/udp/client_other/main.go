package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	port = ":12000"
)

//这里需要注意udp和tcp的不同，udp不需要先开启服务端，DialUDP连接不会报错，而tcp在服务端没开启会直接异常
//所以对于read数据来说，都有阻塞，而write来说，如果写入时，服务端还未开启，那服务端就无法接受到这个数据（该数据就丢失了）
func main() {
	udpAdress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAdress)
	checkError(err)
	log.Println("start write")
	for {
		_, err = conn.Write([]byte("client two"))
		checkError(err)
		// res := []byte{}
		var res [512]byte
		log.Println("start read")
		_, err = conn.Read(res[0:])
		checkError(err)
		log.Println("get mes:", string(res[0:]))
		time.Sleep(time.Second * 5)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
