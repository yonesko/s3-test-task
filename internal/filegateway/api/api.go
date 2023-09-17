package api

import (
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filegateway/service"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
	"log"
	"net/http"
)

func SaveFile(fileGatewayService service.FileGateway) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "POST" {
			http.Error(writer, fmt.Sprintf("method %s is not supported", request.Method), http.StatusBadRequest)
			return
		}
		file, fileHeader, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = fileGatewayService.SaveFile(request.Context(), model.File{
			Name: fileHeader.Filename,
			Body: file,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}

func GetFile(fileGatewayService service.FileGateway) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != "GET" {
			http.Error(writer, fmt.Sprintf("method %s is not supported", request.Method), http.StatusBadRequest)
			return
		}

		file, err := fileGatewayService.GetFile(request.Context(), request.URL.Query().Get("name"))
		//TODO file closer
		if err != nil {
			log.Default().Printf("GetFile err: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(writer, file)
		if err != nil {
			log.Default().Printf("Copy err: %s", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}
