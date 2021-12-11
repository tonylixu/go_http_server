export tag=v1.12
root:
    export ROOT=github.com/tonylixu/go_http_server

build:
	echo "Building http_server binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "Building http_server container"
	docker build -t tonylixu/go_http_server:${tag} .

push: release
	echo "Pushing tonylixu/go_http_server"
	docker push tonylixu/go_http_server:${tag}
