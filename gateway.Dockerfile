FROM golang
LABEL authors="gdanichev"
WORKDIR /app
COPY . ./
RUN GOOS=linux go build -o gateway cmd/gateway/main.go
CMD ["./gateway"]
