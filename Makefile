XFLAGS = -X main.version=$(shell git describe --tags)

dev:
	@ go run cmd/amoeba/main.go

bin:
	@ go build -o bin/amoeba -ldflags "-w ${XFLAGS}" cmd/amoeba/main.go

test:
	@ go test -v ./...

for-osx:
	@ env GOOS=darwin GOARCH=amd64 go build -o bin/amoeba.osx-amd64 -ldflags "-w ${XFLAGS}" cmd/amoeba/main.go

for-linux:
	@ env GOOS=linux GOARCH=amd64 go build -o bin/amoeba.linux-amd64 -ldflags "-w ${XFLAGS}" cmd/amoeba/main.go

for-win:
	@ env GOOS=windows GOARCH=amd64 go build -o bin/amoeba.windows-amd64 -ldflags "-w ${XFLAGS}" cmd/amoeba/main.go


.PHONY: dev bin for-osx for-linux for-win
