package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filestorage/api"
	"github.com/yonesko/s3-test-task/internal/filestorage/memoryfilestorage"
	"github.com/yonesko/s3-test-task/pkg/httplog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var gatewayAddr = flag.String("gateway", "localhost:8361", "gateway to register")

func main() {
	flag.Parse()
	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer stop()

	//TODO move to client, retry
	resp, err := http.Get(fmt.Sprintf("http://%s/register?host=localhost:8208", *gatewayAddr))
	if err != nil {
		log.Fatal("cant register:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("cant register:", resp.Status)
	}

	server := http.Server{Addr: ":8208", ReadTimeout: time.Second * 3}

	fileStorage := memoryfilestorage.NewStorage()
	http.HandleFunc("/file", httplog.Log(func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case "GET":
			api.GetFile(fileStorage)(writer, request)
		case "POST":
			api.SaveFile(fileStorage)(writer, request)
		default:
			http.Error(writer, fmt.Sprintf("method %s is not supported", request.Method), http.StatusBadRequest)
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
