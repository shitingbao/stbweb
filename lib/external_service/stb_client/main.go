package main

import (
	"context"
	"io"
	"log"
	"os"
	"stbweb/lib/external_service/stbserver"
	"strconv"
	"time"

	"github.com/pborman/uuid"

	"google.golang.org/grpc"

	_ "google.golang.org/grpc/balancer/grpclb"
)

const port = "localhost:5000"

func main() {
	startConnect()

}

func startConnect() {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := stbserver.NewStbServerClient(conn) //新建client

	// getSummoner(c)
	// getAllSummoner(c)
	// putSummoner(c)
	// shareSummoner(c)
	// sendfile(c)
	// sendBigFile(c)
	sendGroupFile(c)
}

//普通数据传输
func getSummoner(c stbserver.StbServerClient) {
	character, err := c.GetSummonerInfo(context.Background(), &stbserver.Identity{
		Idcard: "qwer",
		Name:   "shitingbao",
	})
	if err != nil {
		log.Println("err:", err)
	}
	log.Println("character:", character)
}

//单向流，接受值
func getAllSummoner(c stbserver.StbServerClient) {
	req, err := c.GetAllSummonerInfo(context.Background(), &stbserver.Identity{
		Idcard: "qwer",
		Name:   "shitingbao",
	})
	if err != nil {
		log.Println("err:", err)
		return
	}
	for {
		da, err := req.Recv()
		if err != nil {
			log.Println("err:", err)
			break
		}
		log.Println("da:", da)
	}
}

//单向流，发送值
func putSummoner(c stbserver.StbServerClient) {
	res, err := c.PutSummonerInfo(context.Background())
	if err != nil {
		log.Println("err:", err)
		return
	}
	i := 0
	for {
		if i > 2 {
			break
		}
		//这里注意发送后，如果服务端没有接受就关闭了连接，是无法接收到数据的，所以这里加一个timie.Sleep
		if err := res.Send(&stbserver.Identity{
			Idcard: uuid.NewUUID().String(),
			Name:   "shitingbao",
		}); err != nil {
			log.Println("err:", err)
			break
		}
		i++
		time.Sleep(time.Second * 1)
	}
}

//双向流
func shareSummoner(c stbserver.StbServerClient) {
	cli, err := c.ShareSummonerInfo(context.Background())
	if err != nil {
		log.Println("err:", err)
		return
	}
	go func() {
		for {
			da, err := cli.Recv()
			if err != nil {
				log.Println("err:", err)
				return
			}
			log.Println("da:", da)
		}
	}()

	go func() {
		i := 0
		for {
			if i > 3 {
				break
			}
			log.Println("send:", i)
			if err := cli.Send(&stbserver.Identity{
				Idcard: uuid.NewUUID().String(),
				Name:   "shitingbao",
			}); err != nil {
				log.Println("err:", err)
				return
			}
			time.Sleep(time.Second)
			i++
		}

	}()
	time.Sleep(time.Second * 10)
}

func sendfile(c stbserver.StbServerClient) {
	res, err := c.SendFile(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	f, err := os.Open("./test.json")
	if err != nil {
		panic(err)
	}
	sta, err := f.Stat()
	if err != nil {
		panic(err)
	}
	log.Println("size:", sta.Size())
	defer f.Close()
	buf := make([]byte, sta.Size())
	i := 1
	for {
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if err == io.EOF {
			log.Println(err)
			break
		}

		res.Send(&stbserver.FileMessage{
			FileName: strconv.Itoa(i),
			FileType: "json",
			FileData: buf,
			IsCarry:  true,
		})
		i++
	}
	time.Sleep(time.Second * 2)
}

func sendBigFile(c stbserver.StbServerClient) {
	f, err := os.Open("./test.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	// log.Println(fInfo.Size())
	fSize := fInfo.Size()
	i := 1
	res, err := c.SendFile(context.Background())
	if err != nil {
		panic(err)
	}

	for {
		bufSize := 200
		if int64(200*i) > fSize && int64(200*(i-1)) < fSize {
			bufSize = int(fSize) - ((i - 1) * 200)
		}

		buf := make([]byte, bufSize)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if err == io.EOF {
			log.Println(err)
			break
		}
		res.Send(&stbserver.FileMessage{
			FileName: strconv.Itoa(i),
			FileType: "json",
			FileData: buf,
			IsCarry:  true,
		})
		i++
	}
	time.Sleep(time.Second * 2)
	return
}

func sendGroupFile(c stbserver.StbServerClient) {
	res, err := c.SendGroupFile(context.Background())
	if err != nil {
		return
	}
	filename := "test.json"
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	fInfo, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fSize := fInfo.Size()
	i := 1
	for {
		isCarry := false
		isStart := false
		if i == 1 {
			isStart = true
		}
		bufSize := 200
		if int64(200*i) > fSize && int64(200*(i-1)) < fSize {
			bufSize = int(fSize) - ((i - 1) * 200)
			isCarry = true
		}
		buf := make([]byte, bufSize)
		_, err := f.Read(buf)
		if err != nil && err != io.EOF {
			break
		}
		if err == io.EOF {
			log.Println(err)
			break
		}
		res.Send(&stbserver.FileMessage{
			FileName:  filename,
			FileType:  strconv.Itoa(i),
			FileData:  buf,
			IsCarry:   isCarry,
			IsStart:   isStart,
			User:      "shitingbao",
			TotalSize: fSize,
		})
		i++
	}
	time.Sleep(time.Second * 2)
}
