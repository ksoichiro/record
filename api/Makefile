.PHONY: lint
lint:
	GO111MODULE=on golangci-lint run ./...

.PHONY: test
test:
	@GO111MODULE=on go test -covermode=count -coverprofile=c.out ./... && go tool cover -html=c.out -o coverage.html
	@go tool cover -func=c.out | grep total | awk '{ print $$3 }'

.PHONY: install
install:
	GO111MODULE=on go get ./...
