
test:
	go test ./...

lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .
