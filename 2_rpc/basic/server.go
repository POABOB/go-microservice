package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"base/service"
)

func main() {
	// 實例化
	stringService := new(service.StringService)

	// rpc 註冊
	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("Register error: ", registerError)
	}

	// HTTP 協議
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	// 開啟服務
	http.Serve(l, nil)
}
