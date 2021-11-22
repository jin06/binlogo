# todo

# build server mac: go build -o ./bin/binlogo-v1.0.10-darwin-amd64  cmd/server/binlogo.go
#
# version 1.0.11
# mac:
# CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/1.0.11/binlogo-v1.0.11-darwin-amd64/binlogo  cmd/server/binlogo.go
# windows:
# CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o ./bin/1.0.11/binlogo-v1.0.11-windows-amd64/binlogo  cmd/server/binlogo.go
# linux:
# CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ./bin/1.0.11/binlogo-v1.0.11-linux-amd64/binlogo  cmd/server/binlogo.go

# cp -fr configs bin/1.0.11/binlogo-v1.0.11-darwin-amd64
# cp -fr configs bin/1.0.11/binlogo-v1.0.11-windows-amd64
# cp -fr configs bin/1.0.11/binlogo-v1.0.11-linux-amd64

# compress files
#  zip -q -r -o bin/1.0.11/binlogo-v1.0.11-darwin-amd64.zip bin/1.0.11/binlogo-v1.0.11-darwin-amd64
#  zip -q -r -o bin/1.0.11/binlogo-v1.0.11-windows-amd64.zip bin/1.0.11/binlogo-v1.0.11-windows-amd64
#  tar -zcvf bin/1.0.11/binlogo-v1.0.11-linux-amd64.tar.gz bin/1.0.11/binlogo-v1.0.11-linux-amd64


# docker
# docker build -t jin06/binlogo .
# docker
# docker push jin06/binlogo
# docker tag jin06/binlogo jin06/binlogo:1.0.11
# docker push jin06/binlogo:1.0.11

