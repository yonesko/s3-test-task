package filestorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
	"mime/multipart"
	"net/http"
)

type Client interface {
	SaveFile(ctx context.Context, serverUrl string, file model.File) error
	GetFile(ctx context.Context, serverUrl string, name string) (io.Reader, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

// TODO use ctx
func (c client) SaveFile(ctx context.Context, serverUrl string, file model.File) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	part, err := writer.CreateFormFile("file", file.Name)
	if err != nil {
		return fmt.Errorf("CreateFormFile err: %w", err)
	}
	_, err = io.Copy(part, file.Body)
	if err != nil {
		return fmt.Errorf("Copy err: %w", err)
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("http://%s/file", serverUrl), body)
	if err != nil {
		return fmt.Errorf("NewRequest err: %w", err)
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("Do err: %w", err)
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("bas status %s", response.Status)
	}
	return nil
}

func (c client) GetFile(ctx context.Context, serverUrl string, name string) (io.Reader, error) {
	//TODO use custom http client and connection pool
	response, err := http.Get(fmt.Sprintf("http://%s/file?name=%s", serverUrl, name))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("bas status %s", response.Status)
	}
	return response.Body, nil
}
