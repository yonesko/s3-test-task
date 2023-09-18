docker compose up --build

curl -v -F file=@docker-compose.yaml localhost:8361/file &&
  curl "localhost:8361/file?name=docker-compose.yaml" >docker-compose.yaml.d &&
  md5sum docker-compose.yaml*
