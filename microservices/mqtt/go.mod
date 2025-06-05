module github.com/tinh-tinh/tinhtinh/microservices/mqtt

go 1.23.0

toolchain go1.24.1

require (
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/stretchr/testify v1.9.0
	github.com/tinh-tinh/tinhtinh/v2 v2.1.2
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../
