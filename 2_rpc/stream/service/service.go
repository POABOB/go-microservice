package service

import (
	"context"
	"errors"
	"io"
	"log"
	stream_pb "stream/pb"
	"strings"
)

// Service 常數
const (
	StrMaxSize = 1024
)

// Service errors
var (
	ErrMaxSize  = errors.New("超出最大長度 1024 bytes")
	ErrStrValue = errors.New("超出最大長度 1024 bytes")
)

// Class StringService
type StringService struct {
}

// 連接兩字串
func (s *StringService) Concat(ctx context.Context, req *stream_pb.StringRequest) (*stream_pb.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := stream_pb.StringResponse{Ret: ""}
		return &response, nil
	}
	response := stream_pb.StringResponse{Ret: req.A + req.B}
	return &response, nil
}

// 判斷字串相同處
func (s *StringService) Diff(ctx context.Context, req *stream_pb.StringRequest) (*stream_pb.StringResponse, error) {
	if len(req.A) < 1 || len(req.B) < 1 {
		response := stream_pb.StringResponse{Ret: ""}
		return &response, nil
	}
	res := ""
	if len(req.A) >= len(req.B) {
		for _, char := range req.B {
			if strings.Contains(req.A, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range req.A {
			if strings.Contains(req.B, string(char)) {
				res = res + string(char)
			}
		}
	}
	response := stream_pb.StringResponse{Ret: res}
	return &response, nil
}

// server 是 stream
func (s *StringService) LotsOfServerStream(req *stream_pb.StringRequest, qs stream_pb.StringService_LotsOfServerStreamServer) error {
	response := stream_pb.StringResponse{Ret: req.A + req.B}
	for i := 0; i < 10; i++ {
		// 回傳
		qs.Send(&response)
	}
	return nil
}

// client 是stream
func (s *StringService) LotsOfClientStream(qs stream_pb.StringService_LotsOfClientStreamServer) error {
	var params []string
	for {
		// 接收
		in, err := qs.Recv()
		// 接收完全後，回傳順便關閉
		if err == io.EOF {
			qs.SendAndClose(&stream_pb.StringResponse{Ret: strings.Join(params, "")})
			return nil
		}
		if err != nil {
			log.Printf("failed to recv: %v", err)
			return err
		}
		params = append(params, in.A, in.B)
	}
}

// 雙向 stream
func (s *StringService) LotsOfServerAndClientStream(qs stream_pb.StringService_LotsOfServerAndClientStreamServer) error {
	for {
		// 接收
		in, err := qs.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("failed to recv %v", err)
			return err
		}
		// 回傳
		qs.Send(&stream_pb.StringResponse{Ret: in.A + in.B})
	}
}
