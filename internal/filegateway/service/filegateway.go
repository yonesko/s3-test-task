package service

import (
	"context"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
)

type FileGateway interface {
	SaveFile(ctx context.Context, file model.File) error
	GetFile(ctx context.Context, name string) (io.ReadCloser, error)
}
