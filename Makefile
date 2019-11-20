test:
	GO111MODULE=on go test -covermode=count -coverprofile=c.out ./... && go tool cover -html=c.out -o coverage.html

install:
	GO111MODULE=on go get ./...
