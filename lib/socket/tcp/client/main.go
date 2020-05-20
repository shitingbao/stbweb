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

//write和read都是阻塞形式的读写实际连接看实际应用
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	for {
		result := make([]byte, 256)
		_, err = conn.Read(result)
		checkError(err)
		fmt.Println("get mes:", string(result))

		wlen, err := conn.Write([]byte("a1好"))
		checkError(err)
		log.Println("wlen:", wlen)
		// result, err := ioutil.ReadAll(conn)

		time.Sleep(time.Second * 5)
	}
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
