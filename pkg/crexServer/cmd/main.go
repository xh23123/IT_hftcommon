package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	proxy "github.com/xh23123/IT_hftcommon/pkg/crexServer"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	server := proxy.NewServer(8090)
	go func() {
		for range c {
			server.Stop()
		}
	}()
	log.Println(server.Run())
}
