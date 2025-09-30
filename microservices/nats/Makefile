.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: format
format:
	go run mvdan.cc/gofumpt@latest -w -l .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy -v

.PHONY: test
test:
	go test -cover ./...

.PHONY: coverage
coverage:
	go clean -testcache
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html

.PHONY: benchmark
benchmark:
	go test ./... -benchmem -bench=. -run=^Benchmark_$
