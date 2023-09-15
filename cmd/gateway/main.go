package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	server := http.Server{Addr: ":8000", ReadTimeout: time.Second * 3}

	http.HandleFunc("/file", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			return
		}
	})

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("ListenAndServe err:", err)
		}
	}()

	select {
	case <-mainCtx.Done():
		err := server.Shutdown(context.Background())
		if err != nil {
			fmt.Println("Shutdown err:", err)
		}
		return
	}
}
