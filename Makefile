.PHONY: build
app = github.com/jin06/binlogo/v2
version = $version
compileTime = $(shell date)
goVersion = $(shell go version)
output ?= bin
darwinName = binlogo-v$(version)-darwin-amd64
darwinDir = $(output)/$(darwinName)
darwinArmName = binlogo-v$(version)-darwin-arm64
darwinArmDir = $(output)/$(darwinArmName)
windowsName = binlogo-v$(version)-windows-amd64
windowsDir = $(output)/$(windowsName)
linuxName = binlogo-v$(version)-linux-amd64
linuxDir = $(output)/$(linuxName)
armName = binlogo-v$(version)-linux-arm64
armDir = $(output)/$(armName)
buildArgs = -ldflags="-X '$(app)/configs.Version=$(version)' -X '$(app)/configs.BuildTime=$(compileTime)' -X '$(app)/configs.GoVersion=$(goVersion)'" cmd/server/binlogo.go
build:
	mkdir -p $(darwinDir)/etc
	mkdir -p $(darwinArmDir)/etc
	mkdir -p $(windowsDir)/etc
	mkdir -p $(linuxDir)/etc
	mkdir -p $(armDir)/etc
	cp -rf assets $(darwinDir)
	cp -rf assets $(darwinArmDir)
	cp -rf assets $(windowsDir)
	cp -rf assets $(linuxDir)
	cp -rf assets $(armDir)
	cp etc/binlogo.yaml $(darwinDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(darwinArmDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(windowsDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(linuxDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(armDir)/etc/binlogo.yaml
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o $(darwinDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=darwin GOARCH=arm64 go build -o $(darwinArmDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o $(windowsDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o $(linuxDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=arm64 go build -o $(armDir)/binlogo $(buildArgs)
	zip -q -r -o $(output)/$(darwinName).zip  $(output)/$(darwinName)/
	zip -q -r -o $(output)/$(darwinArmName).zip  $(output)/$(darwinArmName)/
	zip -q -r -o $(output)/$(windowsName).zip $(output)/$(windowsName)/
	tar -zcvf $(output/)$(linuxName).tar.gz $(output)/$(linuxName)/
	tar -zcvf $(output/)$(armName).tar.gz $(output)/$(armName)/
.PHONY: docker
version = $version
docker:
	docker build -t jin06/binlogo . --build-arg version=$(version)
	docker push jin06/binlogo
	docker tag jin06/binlogo jin06/binlogo:$(version)
	docker push jin06/binlogo:$(version)

.ONESHELL:

dash:
	sh -c "cd dashboard;NODE_OPTIONS="--openssl-legacy-provider" npm run build:prod"



