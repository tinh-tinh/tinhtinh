module github.com/tinh-tinh/tinhtinh/microservices/redis

go 1.22.2

require (
	github.com/redis/go-redis/v9 v9.7.3
	github.com/stretchr/testify v1.9.0
	github.com/tinh-tinh/tinhtinh/v2 v2.0.0-beta.2
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../
