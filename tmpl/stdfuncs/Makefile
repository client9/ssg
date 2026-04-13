
test:
	go test ./...

lint:
	go mod tidy
	golangci-lint fmt
	golangci-lint run .
env:
	brew upgrade
	brew install golangci-lint
