FROM golang:1.18-alpine as builder

LABEL maintainer="jinlog<jinlong4696@163.com>"

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /binlogo

COPY . .

RUN go build ./cmd/server/binlogo.go

FROM alpine:3.10 as final

#ENV ETCD_ENDPOINTS="127.0.0.1:2379"
#ENV ETCD_PASSWORD=""
#ENV ETCD_USERNAME=""
#ENV NODE_NAME=""
#ENV BINLOGO_ENV="production"
#ENV CONSOLE_LISTEN="0.0.0.0"
#ENV CONSOLE_PORT="9999"
#ENV CLUSTER_NAME="cluster"

COPY --from=builder /binlogo/binlogo /binlogo/binlogo
COPY --from=builder /binlogo/configs/binlogo_docker.yaml /binlogo/configs/binlogo.yaml
COPY --from=builder /binlogo/assets /binlogo/assets
WORKDIR /binlogo

EXPOSE 9999
RUN cd /binlogo

CMD ["/binlogo/binlogo","server", "--config", "/binlogo/configs/binlogo.yaml"]
