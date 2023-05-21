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

// 請求結構體
type StringRequest struct {
	A string
	B string
}

// Service Interface
type Service interface {
	Concat(req StringRequest, ret *string) error
	Diff(req StringRequest, ret *string) error
}

// Class StringService
type StringService struct {
}

// 連接兩字串
func (s StringService) Concat(req StringRequest, ret *string) error {
	// test for length overflow
	if len(req.A)+len(req.B) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	*ret = req.A + req.B
	return nil
}

// 判斷字串相同處
func (s StringService) Diff(req StringRequest, ret *string) error {
	if len(req.A) < 1 || len(req.B) < 1 {
		*ret = ""
		return nil
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
	*ret = res
	return nil
}

// ServiceMiddleware 注入 log 的記錄行為
type ServiceMiddleware func(Service) Service
