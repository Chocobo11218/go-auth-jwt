## run: execute main go application in local
.PHONY: run
run:
	APPENV=local go run app/cmd/main.go

## tidy: special go mod tidy without golang database checksum(GOSUMDB)
.PHONY: tidy
tidy:
	GOSUMDB=off go mod tidy

## test: run go test
test:
	go clean -testcache
	go test -race -v ./...

## test: run go test cover
cover:
	go clean -testcache
	go test -cover ./...
