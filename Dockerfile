FROM golang:1.17.1-alpine as builder
ENV MY_SERVICE_PORT=80
LABEL multi.label1="http"
ADD bin/amd64/go_http_server /go_http_server
EXPOSE 8080

# Add Tini
ENV TINI_VERSION v0.19.0
ADD https://github.com/krallin/tini/releases/download/${TINI_VERSION}/tini /tini
RUN chmod +x /tini
ENTRYPOINT ["/tini", "--"]
CMD ["/go_http_server"]
