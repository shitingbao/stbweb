package main

import (
	"context"
	"log"
	"stbweb/lib/external_service/stbserver"
	"time"

	"github.com/pborman/uuid"

	"google.golang.org/grpc"

	_ "google.golang.org/grpc/balancer/grpclb"
)

const port = "localhost:5000"

func main() {
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := stbserver.NewStbServerClient(conn) //新建client

	// getSummoner(c)
	// getAllSummoner(c)
	putSummoner(c)
	// shareSummoner(c)
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
		time.Sleep(time.Second * 3)
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
