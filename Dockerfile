FROM golang:1.17.1-alpine as builder
ENV MY_SERVICE_PORT=80
LABEL multi.label1="http"
ADD bin/amd64/go_http_server /go_http_server
EXPOSE 8080
ENTRYPOINT ["/go_http_server"]
