package main

import (
	"context"
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filestorage/api"
	"github.com/yonesko/s3-test-task/internal/filestorage/memoryfilestorage"
	"github.com/yonesko/s3-test-task/pkg/httplog"
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

	fileStorage := memoryfilestorage.NewStorage()
	http.HandleFunc("/file", httplog.Log(api.SaveFile(fileStorage)))

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
