.PHONY: build
app = github.com/jin06/binlogo
version = $version
compileTime = $(shell date)
goVersion = $(shell go version)
baseOutput = upload
output = $(baseOutput)/$(version)
darwinDir = $(output)/binlogo-v$(version)-darwin-amd64
windowsDir = $(output)/binlogo-v$(version)-windows-amd64
linuxDir = $(output)/binlogo-v$(version)-linux-amd64
buildArgs = -ldflags="-X '$(app)/configs.Version=$(version)' -X '$(app)/configs.BuildTime=$(compileTime)' -X '$(app)/configs.GoVersion=$(goVersion)'" cmd/server/binlogo.go
build:
	mkdir -p $(darwinDir)/env
	cp env/binlogo.yaml $(darwinDir)/env/binlogo.yaml
	mkdir -p $(windowsDir)/env
	cp env/binlogo.yaml $(windowsDir)/env/binlogo.yaml
	mkdir -p $(linuxDir)/env
	cp env/binlogo.yaml $(linuxDir)/env/binlogo.yaml
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o $(darwinDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o $(windowsDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o $(linuxDir)/binlogo $(buildArgs)
	zip -q -r -o $(baseOutput)/binlogo-darwin-amd64.zip $(darwinDir)
	zip -q -r -o $(baseOutput)/binlogo-windows-amd64.zip $(windowsDir)
	tar -zcvf $(baseOutput)/binlogo-linux-amd64.tar.gz $(linuxDir)

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



