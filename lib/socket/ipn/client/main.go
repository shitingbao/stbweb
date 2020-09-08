package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	service         = ":1200"
	remotelyService = "124.70.156.31:1200"
)

type dataModel struct {
	Adress string
	User   string
	Data   string
}

var (
	dialogList   map[string]net.Conn
	contrastList map[net.Conn]net.Conn //保存两个连接的对应关系，确保连接和断开连接能找到另一个转接连接对象
)

//main监听本地请求,本身是一个server端的服务，在一个连接后连接远程进行转接
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError("tcp:", err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError("listen:", err)
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
	wg := sync.WaitGroup{}
	wg.Add(2)
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	go readAndWrite(&wg, conn)
	// go readAndWrite(&wg, remotelyCon, conn)
	wg.Wait()
}

//连接远程端口
func remotelyConnect(address string) (net.Conn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", remotelyService)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//将一端的读写转发到另一端
func readAndWrite(wg *sync.WaitGroup, con net.Conn) {
	for {
		request := make([]byte, 128)
		readLen, err := con.Read(request)
		if err != nil {
			close(contrastList[con], con, "connect err:"+err.Error())
			break
		}
		if readLen == 0 {
			close(contrastList[con], con, "connect out")
			break
		}
		da := dataModel{}
		if err := json.Unmarshal(request[:readLen], &da); err != nil {
			close(contrastList[con], con, "json err:"+err.Error())
			break
		}
		remotelyCon := dialogList[da.Adress]
		if remotelyCon == nil {
			remotelyCon, err = remotelyConnect(da.Adress)
			if err != nil {
				con.Write([]byte(err.Error()))
				continue
			}
			dialogList[da.Adress] = remotelyCon
			contrastList[con] = remotelyCon
		}
		remotelyCon.Write(request[:readLen])
	}
	wg.Done()
}

func close(remotelyCon, conn net.Conn, reason string) {
	if remotelyCon != nil {
		remotelyCon.Close()
	}
	conn.Close()
	logrus.Info(reason)
}

func checkError(reson string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, reson+": error: %s", err.Error())
		os.Exit(1)
	}
}
