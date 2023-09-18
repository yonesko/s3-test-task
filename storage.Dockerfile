FROM golang
LABEL authors="gdanichev"
WORKDIR /app
COPY . ./
RUN GOOS=linux go build -o storage cmd/storage/main.go
CMD ["./storage", "gateway", "$GATEWAY_ADDR"]
