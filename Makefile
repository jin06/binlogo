.PHONY: build
app = github.com/jin06/binlogo/v2
version = $version
compileTime = $(shell date)
goVersion = $(shell go version)
output ?= bin
darwinName = binlogo-v$(version)-darwin-amd64
darwinDir = $(output)/$(darwinName)
windowsName = binlogo-v$(version)-windows-amd64
windowsDir = $(output)/$(windowsName)
linuxName = binlogo-v$(version)-linux-amd64
linuxDir = $(output)/$(linuxName)
buildArgs = -ldflags="-X '$(app)/configs.Version=$(version)' -X '$(app)/configs.BuildTime=$(compileTime)' -X '$(app)/configs.GoVersion=$(goVersion)'" cmd/server/binlogo.go
build:
	mkdir -p $(darwinDir)/etc
	mkdir -p $(windowsDir)/etc
	mkdir -p $(linuxDir)/etc
	cp -rf assets $(darwinDir)
	cp -rf assets $(windowsDir)
	cp -rf assets $(linuxDir)
	cp etc/binlogo.yaml $(darwinDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(windowsDir)/etc/binlogo.yaml
	cp etc/binlogo.yaml $(linuxDir)/etc/binlogo.yaml
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o $(darwinDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o $(windowsDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o $(linuxDir)/binlogo $(buildArgs)
	zip -q -r -o $(output)/$(darwinName).zip  $(output)/$(darwinName)/
	zip -q -r -o $(output)/$(windowsName).zip $(output)/$(windowsName)/
	tar -zcvf $(output/)$(linuxName).tar.gz $(output)/$(linuxName)/
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
	sh -c "rm -fr assets/dist;mv dashboard/dist ./assets"



