package filestorage

import (
	"context"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
)

type Client interface {
	SaveFile(ctx context.Context, serverUrl string, file model.File) error
	GetFile(ctx context.Context, serverUrl string, name string) (io.Reader, error)
}

type client struct {
}

func (c client) SaveFile(ctx context.Context, serverUrl string, file model.File) error {
	//TODO implement me
	panic("implement me")
}

func (c client) GetFile(ctx context.Context, serverUrl string, name string) (io.Reader, error) {
	//TODO implement me
	panic("implement me")
}
