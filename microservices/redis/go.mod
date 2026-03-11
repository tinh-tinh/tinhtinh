module github.com/tinh-tinh/tinhtinh/microservices/redis

go 1.24.1

require (
	github.com/redis/go-redis/v9 v9.17.3
	github.com/stretchr/testify v1.10.0
	github.com/tinh-tinh/tinhtinh/microservices v1.5.0
	github.com/tinh-tinh/tinhtinh/v2 v2.5.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../

replace github.com/tinh-tinh/tinhtinh/microservices => ../
