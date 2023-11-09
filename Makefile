# todo


# kubectl create -f ./docs/kubernetes/
# build server mac: go build -o ./bin/binlogo-v1.0.41-darwin-amd64  cmd/server/binlogo.go
#
# version 1.0.41
# mac:
# CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/1.0.41/binlogo-v1.0.41-darwin-amd64/binlogo  cmd/server/binlogo.go
# windows:
# CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o ./bin/1.0.41/binlogo-v1.0.41-windows-amd64/binlogo  cmd/server/binlogo.go
# linux:
# CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o ./bin/1.0.41/binlogo-v1.0.41-linux-amd64/binlogo  cmd/server/binlogo.go

# mkdir -p bin/1.0.41/binlogo-v1.0.41-darwin-amd64/configs && cp configs/binlogo.yaml bin/1.0.41/binlogo-v1.0.41-darwin-amd64/configs/binlogo.yaml

# mkdir -p bin/1.0.41/binlogo-v1.0.41-windows-amd64/configs && cp configs/binlogo.yaml bin/1.0.41/binlogo-v1.0.41-windows-amd64/configs/binlogo.yaml

# mkdir -p bin/1.0.41/binlogo-v1.0.41-linux-amd64/configs && cp configs/binlogo.yaml bin/1.0.41/binlogo-v1.0.41-linux-amd64/configs/binlogo.yaml

# compress files
#  zip -q -r -o bin/1.0.41/binlogo-v1.0.41-darwin-amd64.zip bin/1.0.41/binlogo-v1.0.41-darwin-amd64

#  zip -q -r -o bin/1.0.41/binlogo-v1.0.41-windows-amd64.zip bin/1.0.41/binlogo-v1.0.41-windows-amd64

#  tar -zcvf bin/1.0.41/binlogo-v1.0.41-linux-amd64.tar.gz bin/1.0.41/binlogo-v1.0.41-linux-amd64


# docker
# docker build -t jin06/binlogo .
# docker
# docker push jin06/binlogo
# docker tag jin06/binlogo jin06/binlogo:1.0.41
# docker push jin06/binlogo:1.0.41

.PHONY: build
app = github.com/jin06/binlogo
version = 1.0.99
compileTime = $(shell date)
goVersion = $(shell go version)
output = bin/$(version)
darwinDir = $(output)/binlogo-v$(version)-darwin-amd64
windowsDir = $(output)/binlogo-v$(version)-windows-amd64
linuxDir = $(output)/binlogo-v$(version)-linux-amd64
buildArgs = -ldflags="-X '$(app)/configs.Version=$(version)' -X '$(app)/configs.BuildTime=$(compileTime)' -X '$(app)/configs.GoVersion=$(goVersion)'" cmd/server/binlogo.go
build:
	mkdir -p $(darwinDir)/configs
	cp configs/binlogo.yaml $(darwinDir)/configs/binlogo.yaml
	mkdir -p $(windowsDir)/configs
	cp configs/binlogo.yaml $(windowsDir)/configs/binlogo.yaml
	mkdir -p $(linuxDir)/configs
	cp configs/binlogo.yaml $(linuxDir)/configs/binlogo.yaml
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o $(darwinDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o $(windowsDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o $(linuxDir)/binlogo $(buildArgs)
	zip -q -r -o $(output)/binlogo-v$(version)-darwin-amd64.zip $(darwinDir)
	zip -q -r -o $(output)/binlogo-v$(version)-windows-amd64.zip $(windowsDir)
	tar -zcvf $(output)/binlogo-v$(version)-linux-amd64.tar.gz $(linuxDir)

.PHONY: docker
version = 1.0.99
docker:
	docker build -t jin06/binlogo . --build-arg version=$(version)
	docker push jin06/binlogo
	docker tag jin06/binlogo jin06/binlogo:$(version)
	docker push jin06/binlogo:$(version)

	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o ./ -ldflags="-X 'configs.BuildTime=Thu Nov  9 14:42:13 CST 2023'"   cmd/server/binlogo.go
