.PHONY: test-coverage

test-coverage:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html