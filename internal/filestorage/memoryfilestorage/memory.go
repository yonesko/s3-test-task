package memoryfilestorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/yonesko/s3-test-task/internal/filestorage/service"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
	"sync"
)

type storage struct {
	l sync.Mutex
	m map[string][]byte
}

func NewStorage() service.FileStorage {
	return &storage{m: make(map[string][]byte)}
}

func (s *storage) SaveFile(_ context.Context, file model.File) error {
	s.l.Lock()
	defer s.l.Unlock()

	bytes, err := io.ReadAll(file.Body)
	if err != nil {
		return fmt.Errorf("SaveFile ReadAll err:%w", err)
	}
	s.m[file.Name] = bytes
	return nil
}

func (s *storage) GetFile(_ context.Context, name string) (io.Reader, error) {
	s.l.Lock()
	defer s.l.Unlock()
	buf, ok := s.m[name]
	if !ok {
		return &bytes.Buffer{}, nil
	}
	return bytes.NewBuffer(buf), nil
}
