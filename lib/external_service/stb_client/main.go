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
	"google.golang.org/grpc/metadata"
)

const port = "localhost:5000"

func main() {
	startConnect()

}

// CustomerTokenAuth 拦截中间件
type CustomerTokenAuth struct {
}

func (c *CustomerTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (c *CustomerTokenAuth) RequireTransportSecurity() bool {
	return true
}

// 注意，服务器只能配置一个 UnaryInterceptor和StreamClientInterceptor，
// 否则会报错，客户端也是，虽然不会报错，但是只有最后一个才起作用。
// 如果你想配置多个，可以使用拦截器链，如go-grpc-middleware，或者自己实现。
//
// 客户端拦截器
func Clientinterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method == %s ; req == %v ; rep == %v ; duration == %s ; error == %v\n", method, req, reply, time.Since(start), err)
	return err
}

func startConnect() {
	opts := []grpc.DialOption{}
	//grpc.WithInsecure()这个是一定要添加的，代表开启安全的选项
	opts = append(opts, grpc.WithInsecure())

	// 自定义认证(token)，new(myCredential 的时候，由于我们实现了上述2个接口，因此new的时候，程序会执行我们实现的接口
	opts = append(opts, grpc.WithPerRPCCredentials(new(CustomerTokenAuth)))

	// 加上拦截器
	opts = append(opts, grpc.WithUnaryInterceptor(Clientinterceptor))
	// 还有一种如下StreamInterceptor
	// grpc.StreamInterceptor()
	// 还有tls认证
	// WithTransportCredentials，客户端和服务端大同小异，都适用
	// 	服务端的拦截器
	// UnaryServerInterceptor -- 单向调用的拦截器
	// StreamServerInterceptor -- stream调用的拦截器
	// 客户端的拦截器
	// UnaryClientInterceptor
	// StreamClientInterceptor

	conn, err := grpc.Dial(port, opts...)
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
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
	ctx := metadata.AppendToOutgoingContext(context.Background(), "k1", "v1")
	character, err := c.GetSummonerInfo(ctx, &stbserver.Identity{
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
		//最后一次文件大小可能不满200，引起部分不必要的数据流，这里判断出最后一次，大小用总量减去过去发送的所有bufSize的大小来计算
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
