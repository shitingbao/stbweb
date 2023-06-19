package stboutserver

import (
	"context"
	"log"
	"net"
	"test/external_service/stbserver"
	"time"

	"google.golang.org/grpc"
)

// 注意，服务器只能配置一个 UnaryInterceptor和StreamClientInterceptor，
// 否则会报错，客户端也是，虽然不会报错，但是只有最后一个才起作用。
// 如果你想配置多个，可以使用拦截器链，如go-grpc-middleware，或者自己实现。
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	m, err := handler(ctx, req)
	end := time.Now()
	// 记录请求参数 耗时 错误信息等数据
	log.Println(info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	// log.Println("RPC: %s,req:%v start time: %s, end time: %s, err: %v", info.FullMethod, req, start.Format(time.RFC3339), end.Format(time.RFC3339), err)
	return m, err
}

func ServerLoad() {
	port := ":4399"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	opts := []grpc.ServerOption{}
	// 加上拦截器
	opts = append(opts, grpc.UnaryInterceptor(unaryInterceptor))
	s := grpc.NewServer(opts...)
	stbserver.RegisterStbServerServer(s, &StbServe{})
	log.Println("start listen:", port)
	s.Serve(lis)
}
