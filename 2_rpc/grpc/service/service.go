package service

import (
	"context"
	"errors"
	pb "grpc/pb"
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
func (s *StringService) Concat(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	if len(req.A)+len(req.B) > StrMaxSize {
		response := pb.StringResponse{Ret: ""}
		return &response, nil
	}
	response := pb.StringResponse{Ret: req.A + req.B}
	return &response, nil
}

// 判斷字串相同處
func (s *StringService) Diff(ctx context.Context, req *pb.StringRequest) (*pb.StringResponse, error) {
	if len(req.A) < 1 || len(req.B) < 1 {
		response := pb.StringResponse{Ret: ""}
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
	response := pb.StringResponse{Ret: res}
	return &response, nil
}

// // ServiceMiddleware 注入 log 的記錄行為
// type ServiceMiddleware func(Service) Service
