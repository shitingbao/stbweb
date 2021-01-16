package main

import (
	"log"
	"net"
	"strings"
	"time"
)

var (
	controlWaitList = make(map[string]net.Conn)   //控制方连接列表,标识对应
	inviteWaitList  = make(map[string]net.Conn)   //邀请方连接列表
	contrastList    = make(map[net.Conn]net.Conn) //对应表，两边都对应
	loopWait        = make(map[string]chan bool)  //用于两方等待时的阻塞，key为两边对应标识
	port            = ":1200"
	// ConReadDeadline = 15 * time.Minute
	ConReadDeadline = 5 * time.Second
)

func main() {
	tcpAdree, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		log.Println("ResolveTCPAddr err", err)
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAdree)
	if err != nil {
		log.Println("ListenTCP err", err)
		return
	}
	defer func() {
		log.Println("stop")
	}()
	log.Println("start listen :1200")
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept err", err)
			continue
		}
		go handleClient(con)
	}
}

//获取第一次的标识，用于匹配两个连接，标识内容为 fi/fc:(邀请方/控制方)+对应标识编号
//邀请或者控制，第一次连接后验证是否有对应标识编号的另一个连接，否则就等待
func handleClient(con net.Conn) {
	log.Println("con.RemoteAddr:", con.RemoteAddr().String())
	defer func() {
		if err := recover(); err != nil {
			log.Println("handleClient:", err)
		}
	}()
	con.SetReadDeadline(time.Now().Add(ConReadDeadline))
	defer close(con)
	for {
		request := make([]byte, 128)
		readLine, err := con.Read(request)
		if err != nil {
			log.Println("con.Read:", err)
			break
		}
		if readLine == 0 { //out
			break
		}
		mes := string(request[:readLine])
		log.Println("total into:", mes)
		switch {
		//fi fc都为第一次处理
		case strings.HasPrefix(mes, "fi:"):
			log.Println("fi into", mes)
			firstOpera(mes, con, controlWaitList, inviteWaitList)
		case strings.HasPrefix(mes, "fc:"):
			log.Println("fc into", mes)
			firstOpera(mes, con, inviteWaitList, controlWaitList)
		default:
			operaCon := contrastList[con]
			if operaCon == nil {
				break
			}
			operaCon.Write(request[:readLine])
			con.SetReadDeadline(time.Now().Add(ConReadDeadline))
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
		log.Println(user, ":i start wait")
		<-aisle
		log.Println(user, ":i continue")
		con.Write([]byte("success"))
	} else {
		c[user] = con
		log.Println(user, ":c start wait")
		aisle <- true
		log.Println(user, ":c continue")
	}
}

//关闭连接，以及清除保存的连接对应关系
//其中一个连接关闭时，关闭和他对应的另一个连接
func close(con net.Conn) {
	con.Close()
	conn := contrastList[con]
	if conn != nil {
		if err := conn.Close(); err != nil {
			log.Println("conn:", err)
		}
		delete(contrastList, conn)
	}
	delete(contrastList, con)
	log.Println("this is close")
}
