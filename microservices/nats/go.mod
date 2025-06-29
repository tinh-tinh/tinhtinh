module github.com/tinh-tinh/tinhtinh/microservices/nats

go 1.24.1

require (
	github.com/nats-io/nats.go v1.43.0
	github.com/stretchr/testify v1.10.0
	github.com/tinh-tinh/tinhtinh/microservices v1.0.0
	github.com/tinh-tinh/tinhtinh/v2 v2.2.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../
replace github.com/tinh-tinh/tinhtinh/microservices => ../
