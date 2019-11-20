test:
	go test -covermode=count -coverprofile=c.out ./... && go tool cover -html=c.out -o coverage.html
