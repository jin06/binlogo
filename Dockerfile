FROM golang:1.16-alpine as builder

LABEL maintainer="jinlog<jinlong4696@163.com>"

ENV GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

ADD . .

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
ADD /binlogo .

EXPOSE 9999

CMD ["/binlogo/binlogo server --configs/binlogo_docker.yaml"]
