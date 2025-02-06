module github.com/tinh-tinh/tinhtinh/microservices/amqlib

go 1.22.0

require (
	github.com/rabbitmq/amqp091-go v1.10.0
	github.com/stretchr/testify v1.9.0
	github.com/tinh-tinh/tinhtinh/v2 v2.0.0-beta.2
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../
