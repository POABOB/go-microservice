package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	stream_pb "stream/pb"

	"google.golang.org/grpc"
)

func main() {
	// 服務參數
	serviceAddress := "127.0.0.1:1234"
	// 調用 gRPC
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	// Client 調用
	stringClient := stream_pb.NewStringServiceClient(conn)

	sendServerStreamRequest(stringClient)

	sendClientStreamRequest(stringClient)

	sendClientAndServerStreamRequest(stringClient)
}

// server 是 stream
func sendServerStreamRequest(client stream_pb.StringServiceClient) {
	stringReq := &stream_pb.StringRequest{A: "A", B: "B"}
	stream, _ := client.LotsOfServerStream(context.Background(), stringReq)
	for {
		// 迴圈接收 stream
		item, stream_error := stream.Recv()

		// EOF 結束
		if stream_error == io.EOF {
			break
		}

		// 錯誤產生
		if stream_error != nil {
			log.Printf("failed to recv: %v", stream_error)
		}

		// 接收結果
		fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, item.GetRet())
	}
}

// client 是 stream
func sendClientStreamRequest(client stream_pb.StringServiceClient) {
	fmt.Printf("test sendClientStreamRequest \n")

	stream, err := client.LotsOfClientStream(context.Background())
	for i := 0; i < 10; i++ {
		if err != nil {
			log.Printf("failed to call: %v", err)
			break
		}
		stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("failed to recv: %v", err)
	}
	log.Printf("sendClientStreamRequest ret is : %s", reply.Ret)
}

// 雙向 stream
func sendClientAndServerStreamRequest(client stream_pb.StringServiceClient) {
	fmt.Printf("test sendClientAndServerStreamRequest \n")
	var err error
	stream, err := client.LotsOfServerAndClientStream(context.Background())
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}
	var i int
	for {
		err1 := stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
		if err1 != nil {
			log.Printf("failed to send: %v", err)
			break
		}
		reply, err2 := stream.Recv()
		if err2 != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("sendClientAndServerStreamRequest Ret is : %s", reply.Ret)
		i++

		if i == 100 {
			break
		}
	}
}
