package upgrade

import (
	"log"
	"net"
)

func NewTcpCustomer() error {
	add, err := net.ResolveTCPAddr("tcp4", ":8080")
	if err != nil {
		return err
	}
	listen, err := net.ListenTCP("tcp", add)
	if err != nil {
		return err
	}
	for {
		con, err := listen.Accept()
		if err != nil {
			return err
		}
		go handle(con)
	}
}

func handle(con net.Conn) {
	for {
		r := []byte{}
		if _, err := con.Read(r); err != nil {
			return
		}
		log.Println(string(r))
		con.Write([]byte("hello world"))
	}
}
