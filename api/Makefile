.PHONY: lint
lint:
	GO111MODULE=on golangci-lint run ./...

.PHONY: test
test:
	GO111MODULE=on go test -covermode=count -coverprofile=c.out ./... && go tool cover -html=c.out -o coverage.html

.PHONY: install
install:
	GO111MODULE=on go get ./...
