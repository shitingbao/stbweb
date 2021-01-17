//总体逻辑
//邀请方首先连接至服务端，发送身份标识码，邀请请求进入邀请列表
//控制方连接至服务端，并在等待连接的邀请列表中查找对应标识的连接，匹配成功，连接，匹配失败则断开连接
//要邀请方先连接的逻辑是因为，如果一起连接不好判断map中的标识，容易死锁（map相互找不到），如果用sync.map就需要断言获取con，可能出现问题（因为都是指针和地址）
//注意，任意一方断开，都会使两边同时断开
//注意，要先邀请方连接，控制方再去连接，控制方先连接会直接断开
//注意，两边连接都会反馈一个成功标识，也就是说客户端首次需要发送，并接收一次，类似http的三次握手
package main

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"
)

var (
	contrastList = make(map[string]ConMatchList) //对应表，两边都对应

	port = ":1200"
	// ConReadDeadline = 15 * time.Minute
	ConReadDeadline = 5 * time.Second

	//InviteFlag 邀请者前缀标识
	InviteFlag = "ep-ivt:"

	//ControlFlag 控制方前缀标识
	ControlFlag = "ep-ctl:"
)

//ConMatchList 两方对接成功的连接
type ConMatchList struct {
	Invite, Control net.Conn
}

func main() {
	tcpAdree, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		return
	}
	listener, err := net.ListenTCP("tcp", tcpAdree)
	if err != nil {
		return
	}
	defer func() {
		log.Println("stop")
	}()
	log.Println("start listen :1200")
	for {
		con, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(con)
	}
}

//获取第一次的标识，用于匹配两个连接，标识内容为 (邀请方/控制方)+对应标识编号
//邀请或者控制，第一次连接后验证是否有对应标识编号的另一个连接，否则就等待
func handleClient(con net.Conn) {
	// log.Println("con.RemoteAddr:", con.RemoteAddr().String())
	defer func() {
		if err := recover(); err != nil {
			log.Println("handleClient:", err)
		}
	}()
	user, _, err := firstOpera(con)
	defer close(user)
	if err != nil {
		return
	}
	con.SetReadDeadline(time.Now().Add(ConReadDeadline))
	//控制方连接后，操作一个数据读取即可，邀请方不需要，只需要加入超时时间即可
	client(user, con)
}

func client(user string, con net.Conn) {
	for {
		request := make([]byte, 128)
		readLine, err := con.Read(request)
		if err != nil {
			break
		}
		tp := contrastList[user]
		if tp.Invite == nil {
			break
		}
		if _, err := tp.Invite.Write(request[:readLine]); err != nil {
			return
		}
		(tp.Control).SetReadDeadline(time.Now().Add(ConReadDeadline))
		(tp.Invite).SetReadDeadline(time.Now().Add(ConReadDeadline))
	}
}

//第一次连接时的操作，邀请方进入等待连接列表，或者控制方匹配，反馈成功后的user标识或者err
//不以ep-invite，ep-control开头的都反馈no connect
//role标识是控制方还是邀请方，用于连接后两边约束收发，待定
func firstOpera(con net.Conn) (user, role string, err error) {
	request := make([]byte, 128)
	readLine, err := con.Read(request)
	if err != nil {
		return user, role, err
	}
	if readLine == 0 { //out
		return user, role, errors.New("no connect")
	}
	mes := string(request[:readLine])
	user = mes[7:]
	switch {
	case strings.HasPrefix(mes, InviteFlag):
		contrastList[user] = ConMatchList{Invite: con}
		role = InviteFlag
		con.Write([]byte("ok"))
	case strings.HasPrefix(mes, ControlFlag):
		if !controlOpera(user, con) { //没有对应连接就退出
			return user, role, errors.New("no connect")
		}
		role = ControlFlag
		con.Write([]byte("success"))
	default:
		return user, role, errors.New("no connect")
	}
	return user, role, nil
}

//控制方连接后，找到对应邀请方，连接后，邀请方删除等待列表中的标识
//这里统一把连接放入contrastList中，断开在close中统一处理
//所以先执行contrastList赋值
//注意修改了map里的指，需要放回map，也就是135行，容易混
func controlOpera(user string, con net.Conn) bool {
	conList := contrastList[user]
	conList.Control = con
	contrastList[user] = conList
	if conList.Invite == nil {
		return false
	}
	return true
}

//关闭连接，以及清除保存的连接对应关系
//其中一个连接关闭时，关闭和他对应的另一个连接
func close(user string) {
	tp := contrastList[user]
	if (tp.Control) != nil {
		(tp.Control).Close()
	}
	if (tp.Invite) != nil {
		(tp.Invite).Close()
	}
	delete(contrastList, user)
}
