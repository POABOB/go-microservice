package main

import (
	"context"
	"fmt"
	pb "grpc/pb"

	"google.golang.org/grpc"
)

func main() {
	// 調用 gRPC
	serviceAddress := "127.0.0.1:1234"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	// 最後關閉
	defer conn.Close()

	// 開啟 Client
	bookClient := pb.NewStringServiceClient(conn)
	stringReq := &pb.StringRequest{A: "A", B: "B"}

	// 調用方法
	reply, _ := bookClient.Concat(context.Background(), stringReq)
	fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, reply.Ret)
}
