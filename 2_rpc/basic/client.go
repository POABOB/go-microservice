package main

import (
	"fmt"
	"log"
	"net/rpc"

	"base/service"
)

func main() {
	// HTTP 呼叫
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// 同步呼叫
	stringReq := &service.StringRequest{"A", "B"}
	var reply string
	err = client.Call("StringService.Concat", stringReq, &reply)
	if err != nil {
		log.Fatal("StringService error:", err)
	}
	fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, reply)

	// 異步呼叫
	stringReq = &service.StringRequest{"ACD", "BDF"}
	call := client.Go("StringService.Diff", stringReq, &reply, nil)
	_ = <-call.Done
	fmt.Printf("StringService Diff : %s diff %s = %s\n", stringReq.A, stringReq.B, reply)
}
