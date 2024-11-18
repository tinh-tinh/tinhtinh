.PHONY: test-coverage

test-coverage:
	go clean -testcache
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	rm -rf middleware/logger/logs