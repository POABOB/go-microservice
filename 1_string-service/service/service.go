package service

import (
	"errors"
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

// Service Interface
type Service interface {
	// 連接兩字串
	Concat(a, b string) (string, error)

	// 判斷字串相同處
	Diff(a, b string) (string, error)

	// HealthCheck check service health status
	HealthCheck() bool
}

// Class StringService
type StringService struct {
}

// 連接兩字串
func (s StringService) Concat(a, b string) (string, error) {
	// test for length overflow
	if len(a)+len(b) > StrMaxSize {
		return "", ErrMaxSize
	}
	return a + b, nil
}

// 判斷字串相同處
func (s StringService) Diff(a, b string) (string, error) {
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

// HealthCheck implement Service method
// 只返回true，暫時不實現
func (s StringService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
