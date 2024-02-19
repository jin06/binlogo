FROM golang:1.21.4-alpine3.18 as builder

LABEL maintainer="jinlog<jinlong4696@163.com>"

ENV GO111MODULE=on \
#    GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /binlogo

COPY . .

ARG VERSION
ENV app=github.com/jin06/binlogo

RUN go mod vendor

RUN go build  -ldflags="-X '$app/configs.Version=$VERSION' -X '$app/configs.BuildTime=$(date)' -X '$app/configs.GoVersion=$(go env GOVERSION)'" ./cmd/server/binlogo.go

RUN ./binlogo version

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
COPY --from=builder /binlogo/etc/binlogo_docker.yaml /binlogo/etc/binlogo.yaml
COPY --from=builder /binlogo/assets /binlogo/assets
WORKDIR /binlogo

EXPOSE 9999
RUN cd /binlogo

CMD ["/binlogo/binlogo","server", "--config", "/binlogo/etc/binlogo.yaml"]
