package main

import (
	"log"
	"net"
	"strings"
	"time"
)

type conCoat struct {
	Con  net.Conn
	User string
}

var (
	controlWaitList map[string]net.Conn   //控制方连接列表,标识对应
	inviteWaitList  map[string]net.Conn   //邀请方连接列表
	contrastList    map[net.Conn]net.Conn //对应表，两边都对应
	loopWait        map[string]chan bool  //用于两方等待时的阻塞，key为两边对应标识
)

func main() {
	tcpAdree, err := net.ResolveTCPAddr("tcp4", ":1200")
	if err != nil {
	}
	listener, err := net.ListenTCP("tcp", tcpAdree)
	if err != nil {
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(con)
	}
}

//获取第一次的标识，用于匹配两个连接，标识内容为 first+邀请方/控制方+对应标识编号
//邀请或者控制器，第一次连接后验证是否有对应连接，否则就等待
func handleClient(con net.Conn) {
	con.SetReadDeadline(time.Now().Add(2 * time.Minute))
	// con.SetReadDeadline(time.Now().Add(2 * time.Second))
	defer close(con)
	for {
		request := make([]byte, 128)
		readLine, err := con.Read(request)
		if err != nil {
			break
		}
		if readLine == 0 { //out
			break
		}
		mes := string(request[:readLine])
		switch {
		//fi fc都为第一次处理
		case strings.HasPrefix(mes, "fi:"):
			firstOpera(mes, con, controlWaitList, inviteWaitList)
		case strings.HasPrefix(mes, "fc:"):
			firstOpera(mes, con, inviteWaitList, controlWaitList)
		default:
			operaCon := contrastList[con]
			operaCon.Write(request[:readLine])
		}
	}
}

//第一次连接时的操作，mes时获取的信息，con本地连接,i，c为连接列表
func firstOpera(mes string, con net.Conn, i, c map[string]net.Conn) {
	user := mes[3:]
	aisle := loopWait[user]
	if aisle == nil {
		aisle = make(chan bool)
		loopWait[user] = aisle
	}
	icon := i[user]
	if icon != nil {
		contrastList[icon] = con
		contrastList[con] = icon
		<-aisle
		con.Write([]byte("success"))
	} else {
		c[user] = con
		aisle <- true
	}
}

//关闭连接，以及清除保存的连接对应关系
func close(con net.Conn) {
	con.Close()
	delete(contrastList, con)
	conn := contrastList[con]
	if conn != nil {
		conn.Close()
		delete(contrastList, conn)
	}
	log.Println("this is close")
}