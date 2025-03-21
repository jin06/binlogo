FROM golang:1.23-alpine as builder

LABEL maintainer="kc<jlonmyway@gmail.com>"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /binlogo

COPY . .

ARG VERSION
ENV app=github.com/jin06/binlogo/v2

RUN go build  -ldflags="-X '$app/configs.Version=$VERSION' -X '$app/configs.BuildTime=$(date)' -X '$app/configs.GoVersion=$(go env GOVERSION)'" ./cmd/server/binlogo.go

RUN ./binlogo version

FROM alpine:3.10 as final

#ENV NODE_NAME=""
#ENV BINLOGO_ENV="production"
#ENV CONSOLE_LISTEN="0.0.0.0"
#ENV CONSOLE_PORT="9999"
#ENV CLUSTER_NAME="cluster"
COPY --from=builder /binlogo/binlogo /binlogo/binlogo
COPY --from=builder /binlogo/etc/binlogo_docker.yaml /binlogo/etc/binlogo.yaml
COPY --from=builder /binlogo/static /binlogo/static
WORKDIR /binlogo

EXPOSE 8081
RUN cd /binlogo

CMD ["/binlogo/binlogo","server", "--config", "/binlogo/etc/binlogo.yaml"]
