package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	port = ":12000"
)

//这里需要注意udp和tcp的不同，udp不需要先开启服务端，DialUDP连接不会报错，而tcp在服务端没开启会直接异常
//所以对于read数据来说，都有阻塞，而write来说，如果写入时，服务端还未开启，那服务端就无法接受到这个数据（该数据就丢失了）
func main() {
	udpAdress, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)
	conn, err := net.ListenUDP("udp", udpAdress)
	checkError(err)
	for {
		log.Println("start handle")
		handle(conn)
	}
}

func handle(conn *net.UDPConn) {
	// res := []byte{}
	var buf [512]byte
	log.Println("start read")
	_, addr, err := conn.ReadFromUDP(buf[0:])
	checkError(err)
	log.Println("get mes:", string(buf[0:]))

	log.Println("start write")
	_, err = conn.WriteToUDP([]byte("this is udp server"), addr)
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s", err.Error())
		os.Exit(1)
	}
}
