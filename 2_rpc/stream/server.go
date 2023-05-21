package main

import (
	"flag"
	"net"
	pb "stream/pb"
	"stream/service"

	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

func main() {
	// 服務參數
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// 開啟 gRRPC Server
	grpcServer := grpc.NewServer()
	// 註冊 Service
	pb.RegisterStringServiceServer(grpcServer, new(service.StringService))
	// 服務
	grpcServer.Serve(lis)
}
