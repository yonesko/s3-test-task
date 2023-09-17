package memoryfilegateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/yonesko/s3-test-task/internal/filegateway/client/filestorage"
	"github.com/yonesko/s3-test-task/internal/filegateway/service"
	"github.com/yonesko/s3-test-task/internal/model"
	"golang.org/x/sync/errgroup"
	"io"
	"sync"
)

type fileParts struct {
	storageServers []string //file storage server urls
}

func toPartName(name string, i int) string {
	return fmt.Sprintf("%s_part_%d", name, i)
}

type gateway struct {
	fileStorageClient filestorage.Client
	l                 sync.Mutex //TODO use multiple locks
	files             map[string]fileParts
	servers           map[string]bool
}

func NewGateway() service.FileGateway {
	return &gateway{files: make(map[string]fileParts), servers: make(map[string]bool)}
}

func (s *gateway) RegisterFileStorageServer(url string) error {
	s.l.Lock()
	defer s.l.Unlock()
	s.servers[url] = true
	return nil
}

func (s *gateway) SaveFile(ctx context.Context, file model.File) error {
	s.l.Lock()
	defer s.l.Unlock()
	storageServers := lo.Keys(s.servers)

	buf, err := io.ReadAll(file.Body) //TODO avoid read all, using content size header
	if err != nil {
		return fmt.Errorf("ReadAll err: %w", err)
	}
	group := errgroup.Group{}
	for i, chunk := range lo.Chunk(buf, len(storageServers)) {
		chunk := chunk
		i := i
		group.Go(func() error {
			err = s.fileStorageClient.SaveFile(ctx, storageServers[i], model.File{
				Name: toPartName(file.Name, i),
				Body: bytes.NewBuffer(chunk),
			})
			if err != nil {
				return fmt.Errorf("SaveFile err: %w", err)
			}
			return nil
		})
	}
	err = group.Wait()
	if err != nil {
		return err
	}

	s.files[file.Name] = fileParts{storageServers: storageServers}
	return nil
}

// TODO rewrite to concurrent
func (s *gateway) GetFile(ctx context.Context, name string) (io.Reader, error) {
	s.l.Lock()
	defer s.l.Unlock()
	parts, ok := s.files[name]
	if !ok {
		return &bytes.Buffer{}, nil
	}
	buffer := &bytes.Buffer{}
	for i, storageServer := range parts.storageServers {
		partName := toPartName(name, i)
		file, err := s.fileStorageClient.GetFile(ctx, storageServer, partName)
		if err != nil {
			return &bytes.Buffer{}, nil
		}
		buf, err := io.ReadAll(file)
		if err != nil {
			return &bytes.Buffer{}, fmt.Errorf("ReadAll err: %w", err)
		}
		_, err = buffer.Write(buf)
		if err != nil {
			return &bytes.Buffer{}, fmt.Errorf("Write err: %w", err)
		}
	}

	return buffer, nil
}
