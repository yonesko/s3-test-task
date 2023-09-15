package api

import (
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filestorage/service"
	"github.com/yonesko/s3-test-task/internal/model"
	"net/http"
)

func SaveFile(fileStorageService service.FileStorage) http.HandlerFunc {
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
		err = fileStorageService.SaveFile(request.Context(), model.File{
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
