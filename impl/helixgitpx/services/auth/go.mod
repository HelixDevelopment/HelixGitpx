module github.com/helixgitpx/helixgitpx/services/auth

go 1.23.0

toolchain go1.23.4

require github.com/helixgitpx/platform v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241202173237-19429a94021a // indirect
	google.golang.org/grpc v1.69.2 // indirect
	google.golang.org/protobuf v1.36.0 // indirect
)

replace github.com/helixgitpx/platform => ../../platform

replace github.com/helixgitpx/helixgitpx/gen => ../../gen
