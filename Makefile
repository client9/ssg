
test:
	go test ./...

lint:
	go mod tidy
	gofmt -w -s *.go
	golangci-lint run .
env:
	brew upgrade
	brew install golangci-lint
	go install golang.org/x/tools/cmd/goimports@latest
