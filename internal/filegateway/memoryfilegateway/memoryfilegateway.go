package memoryfilegateway

import (
	"bytes"
	"context"
	"fmt"
	"github.com/samber/lo"
	"github.com/yonesko/s3-test-task/internal/filegateway/client/filestorage"
	"github.com/yonesko/s3-test-task/internal/filegateway/service"
	"github.com/yonesko/s3-test-task/internal/model"
	"io"
	"log"
	"math"
	"sync"
)

var (
	ErrNoStorageAvailable = fmt.Errorf("no storageServers availble")
)

type filePart struct {
	storageServer string
	partNum       int
	originName    string
}

func (p filePart) name() string {
	return fmt.Sprintf("%s_part_%d", p.originName, p.partNum)
}

type gateway struct {
	fileStorageClient filestorage.Client
	l                 sync.Mutex //TODO use multiple locks
	files             map[string][]filePart
	storageServers    []string //available storageServers
}

func NewGateway(fileStorageClient filestorage.Client) service.FileGateway {
	return &gateway{files: make(map[string][]filePart), fileStorageClient: fileStorageClient}
}

func (s *gateway) RegisterFileStorageServer(url string) error {
	s.l.Lock()
	defer s.l.Unlock()
	if !lo.Contains(s.storageServers, url) {
		log.Println(url, "is registered")
		s.storageServers = append(s.storageServers, url)
	}
	return nil
}

// TODO rewrite to concurrent
func (s *gateway) SaveFile(ctx context.Context, file model.File) error {
	s.l.Lock()
	defer s.l.Unlock()
	if len(s.storageServers) == 0 {
		return ErrNoStorageAvailable
	}

	buf, err := io.ReadAll(file.Body) //TODO avoid read all, using content size header
	if err != nil {
		return fmt.Errorf("ReadAll err: %w", err)
	}
	var parts []filePart
	for i, chunk := range lo.Chunk(buf, int(math.Ceil(float64(len(buf))/float64(len(s.storageServers))))) {
		part := filePart{
			storageServer: s.storageServers[i],
			partNum:       i,
			originName:    file.Name,
		}
		err = s.fileStorageClient.SaveFile(ctx, part.storageServer, model.File{
			Name: part.name(),
			Body: bytes.NewBuffer(chunk),
		})
		if err != nil {
			return fmt.Errorf("SaveFile err: %w", err)
		}
		parts = append(parts, part)
	}

	s.files[file.Name] = parts
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
	for _, part := range parts {
		partName := part.name()
		file, err := s.fileStorageClient.GetFile(ctx, part.storageServer, partName)
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
