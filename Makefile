mockgen:
	mockgen -destination  internal/mock/client/filestorage.go -package mock_client \
	github.com/yonesko/s3-test-task/internal/filegateway/client/filestorage Client

test:
	docker compose up