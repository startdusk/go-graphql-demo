mock:
	mockery --all --keeptree

test: clean
	go test ./...

clean:
	go mod tidy
	go vet ./...
	go fmt ./...