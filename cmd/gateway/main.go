package main

import (
	"context"
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filegateway/api"
	"github.com/yonesko/s3-test-task/internal/filegateway/client/filestorage"
	"github.com/yonesko/s3-test-task/internal/filegateway/memoryfilegateway"
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

	server := http.Server{Addr: ":8361", ReadTimeout: time.Second * 3}

	fileGateway := memoryfilegateway.NewGateway(filestorage.NewClient())
	http.HandleFunc("/file", httplog.Log(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			api.GetFile(fileGateway)(writer, request)
		case "POST":
			api.SaveFile(fileGateway)(writer, request)
		default:
			http.Error(writer, fmt.Sprintf("method %s is not supported", request.Method), http.StatusBadRequest)
		}
	}))
	http.HandleFunc("/register", httplog.Log(func(writer http.ResponseWriter, request *http.Request) {
		err := fileGateway.RegisterFileStorageServer(request.URL.Query().Get("host"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}))

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
