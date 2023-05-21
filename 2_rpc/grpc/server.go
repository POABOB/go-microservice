package main

import (
	"net"

	pb "grpc/pb"
	"grpc/service"

	"flag"

	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	// 監聽地址 端口
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 初始化 Server
	grpcServer := grpc.NewServer()

	// 實例化，gRPC 註冊
	pb.RegisterStringServiceServer(grpcServer, new(service.StringService))

	// 開啟服務
	grpcServer.Serve(lis)
}
