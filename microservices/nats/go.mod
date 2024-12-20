module github.com/tinh-tinh/tinhtinh/microservices/nats

go 1.22.2

require (
	github.com/nats-io/nats.go v1.38.0
	github.com/stretchr/testify v1.9.0
	github.com/tinh-tinh/tinhtinh/v2 v2.0.0-beta.2
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../
