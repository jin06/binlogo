.PHONY: build
app = github.com/jin06/binlogo
version = $version
compileTime = $(shell date)
goVersion = $(shell go version)
baseOutput = upload
output = $(baseOutput)
darwinDirName = binlogo-v$(version)-darwin-amd64 
darwinDir = $(output)/$(darwinDirName)
windowsDirName = binlogo-v$(version)-linux-amd64
windowsDir = $(output)/$(windowsDirName)
linuxDirDirName = binlogo-v$(version)-linux-amd64
linuxDir = $(output)/$(linuxDirDirName)
buildArgs = -ldflags="-X '$(app)/configs.Version=$(version)' -X '$(app)/configs.BuildTime=$(compileTime)' -X '$(app)/configs.GoVersion=$(goVersion)'" cmd/server/binlogo.go
build:
	mkdir -p $(darwinDir)/etc
	cp etc/binlogo.yaml $(darwinDir)/etc/binlogo.yaml
	mkdir -p $(windowsDir)/etc
	cp etc/binlogo.yaml $(windowsDir)/etc/binlogo.yaml
	mkdir -p $(linuxDir)/etc
	cp etc/binlogo.yaml $(linuxDir)/etc/binlogo.yaml
	CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o $(darwinDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -o $(windowsDir)/binlogo $(buildArgs)
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o $(linuxDir)/binlogo $(buildArgs)
	cd $(output)
	zip -q -r -o ./binlogo-darwin-amd64.zip $(darwinDir)
	zip -q -r -o ./binlogo-windows-amd64.zip $(windowsDir)
	tar -zcvf ./binlogo-linux-amd64.tar.gz $(linuxDir)
	cd ..
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



