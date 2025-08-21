module github.com/tinh-tinh/tinhtinh/microservices/grpc

go 1.24.1

require (
	github.com/tinh-tinh/tinhtinh/microservices v1.1.0
	github.com/tinh-tinh/tinhtinh/v2 v2.3.1
)

require (
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250818200422-3122310a409c // indirect
	google.golang.org/grpc v1.75.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)

replace github.com/tinh-tinh/tinhtinh/v2 => ../../

replace github.com/tinh-tinh/tinhtinh/microservices => ../
