package memoryfilegateway

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock_client "github.com/yonesko/s3-test-task/internal/mock/client"
	"github.com/yonesko/s3-test-task/internal/model"
	"strings"
	"testing"
)

func TestGateway_SaveFile(t *testing.T) {
	t.Run("single servers save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()

		fileStorageClient := mock_client.NewMockClient(ctrl)
		fileStorageClient.EXPECT().
			SaveFile(ctx, "localhost:6543", gomock.Any()).
			Return(nil)

		fileGateway := NewGateway(fileStorageClient)

		file, err := fileGateway.GetFile(ctx, "dig")
		assert.NoError(t, err)
		assert.Empty(t, file)

		err = fileGateway.SaveFile(ctx, model.File{Name: "dig", Body: strings.NewReader(`{"0Okn", "hx3DzDd", "LeA", "zBnYx9", "T3ON8CAV"}`)})
		assert.ErrorIs(t, err, ErrNoStorageAvailable)

		err = fileGateway.RegisterFileStorageServer("localhost:6543")
		assert.NoError(t, err)

		err = fileGateway.SaveFile(ctx, model.File{Name: "dig", Body: strings.NewReader(`{"0Okn", "hx3DzDd", "LeA", "zBnYx9", "T3ON8CAV"}`)})
		assert.NoError(t, err)
	})
	t.Run("multiple servers save", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx := context.Background()

		fileStorageClient := mock_client.NewMockClient(ctrl)
		fileStorageClient.EXPECT().
			SaveFile(ctx, "localhost:108", gomock.Any()).
			Return(nil)
		fileStorageClient.EXPECT().
			SaveFile(ctx, "localhost:490", gomock.Any()).
			Return(nil)
		fileStorageClient.EXPECT().
			SaveFile(ctx, "localhost:997", gomock.Any()).
			Return(nil)
		fileStorageClient.EXPECT().
			SaveFile(ctx, "localhost:473", gomock.Any()).
			Return(nil)

		fileGateway := NewGateway(fileStorageClient)

		for _, serv := range []string{"localhost:108", "localhost:490", "localhost:997", "localhost:473"} {
			err := fileGateway.RegisterFileStorageServer(serv)
			assert.NoError(t, err)
		}

		err := fileGateway.SaveFile(ctx, model.File{Name: "dig", Body: strings.NewReader(`{"0Okn", "hx3DzDd", "LeA", "zBnYx9", "T3ON8CAV"}`)})
		assert.NoError(t, err)
	})
}
