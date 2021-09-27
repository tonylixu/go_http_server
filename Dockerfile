FROM ubuntu
ENV MY_SERVICE_PORT=80
LABEL multi.label1="http"
ADD bin/amd64/go_http_server /go_http_server
EXPOSE 80
ENTRYPOINT /go_http_server
