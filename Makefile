# todo

# build server mac: go build -o ./bin/binlogo-v1.0.10-darwin-amd64  cmd/server/binlogo.go
#
# mac:
# CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/binlogo-v1.0.10-darwin-amd64  cmd/server/binlogo.go
# windows:
# CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o ./bin/binlogo-v1.0.10-windows-amd64  cmd/server/binlogo.go
# linux:
# CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ./bin/binlogo-v1.0.10-linux-amd64  cmd/server/binlogo.go

# docker


# compress files
#  zip -q -r -o bin/1.0.10/binlogo-v1.0.10-darwin-amd64.zip bin/1.0.10/binlogo-v1.0.10-darwin-amd64
#  zip -q -r -o bin/1.0.10/binlogo-v1.0.10-windows-amd64.zip bin/1.0.10/binlogo-v1.0.10-windows-amd64
#  tar -zcvf bin/1.0.10/binlogo-v1.0.10-linux-amd64.tar.gz bin/1.0.10/binlogo-v1.0.10-linux-amd64
