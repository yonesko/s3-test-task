package memoryfilestorage

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/yonesko/s3-test-task/internal/model"
	"testing"
)

func TestStorage(t *testing.T) {
	t.Run("read non-existing", func(t *testing.T) {
		fileStorage := NewStorage()
		file, err := fileStorage.GetFile(context.Background(), "numerous")
		assert.NoError(t, err)
		assert.Empty(t, file)
	})
	t.Run("write and read", func(t *testing.T) {
		fileStorage := NewStorage()

		err := fileStorage.SaveFile(context.Background(), model.File{
			Name: "numerous",
			Body: bytes.NewBufferString(`lkmjhgfdtrszerxdcftvgyhbjknlm`),
		})
		assert.NoError(t, err)

		file, err := fileStorage.GetFile(context.Background(), "numerous")
		assert.NoError(t, err)
		assert.Equal(t, file, bytes.NewBufferString(`lkmjhgfdtrszerxdcftvgyhbjknlm`))
	})
}
