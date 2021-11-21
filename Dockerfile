FROM golang:1.16-alpine as builder

LABEL maintainer="jinlog<jinlong4696@163.com>"

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ADD . /binlogo

RUN go mod vendor

RUN go build -o /binlogo/binlogo /binlogo/cmd/server/binlogo.go

FROM scratch as final
ENV ETCD_ENDPOINTS="127.0.0.1:2379" \
    ETCD_PASSWORD="" \
    NODE_NAME="" \
    BINLOGO_ENV="production" \
    CONSOLE_LISTEN="0.0.0.0" \
    CONSOLE_PORT="9999" \
    CLUSTER_NAME="cluster"
COPY --from=builder /binlogo /binlogo

EXPOSE 9999

CMD ["/binlogo/binlogo server --/binlogo/configs/binlogo_docker.yaml"]
