package handlers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/POABOB/go-microservice/ticket-system/pb/user/user-service/svc"
)

func InterruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)

	// Place whatever shutdown handling you want here

	errc <- terminateError
}

func SetConfig(cfg svc.Config) svc.Config {
	return cfg
}
