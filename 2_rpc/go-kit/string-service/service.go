package string_service

import (
	"context"
	"errors"
	"fmt"
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

// Service interface
type Service interface {
	Concat(ctx context.Context, a, b string) (string, error)
	Diff(ctx context.Context, a, b string) (string, error)
}

// Class StringService
type StringService struct {
}

// 連接兩字串
func (s StringService) Concat(ctx context.Context, a, b string) (string, error) {
	// test for length overflow
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	fmt.Printf("StringService Concat return %s", a+b)
	return a + b, nil
}

// 判斷字串相同處
func (s StringService) Diff(ctx context.Context, a, b string) (string, error) {
	if len(a) < 1 || len(b) < 1 {
		return "", nil
	}
	res := ""
	if len(a) >= len(b) {
		for _, char := range b {
			if strings.Contains(a, string(char)) {
				res = res + string(char)
			}
		}
	} else {
		for _, char := range a {
			if strings.Contains(b, string(char)) {
				res = res + string(char)
			}
		}
	}
	return res, nil
}

// ServiceMiddleware 注入 log 的記錄行為
type ServiceMiddleware func(Service) Service
