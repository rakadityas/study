module protocols/grpc

go 1.20

require (
	github.com/lib/pq v1.10.9
	google.golang.org/grpc v1.60.1
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/protobuf v1.36.4 // indirect
)

replace google.golang.org/genproto/googleapis/rpc => google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80

replace google.golang.org/protobuf => google.golang.org/protobuf v1.33.0

replace golang.org/x/net => golang.org/x/net v0.20.0

replace golang.org/x/sys => golang.org/x/sys v0.16.0

replace golang.org/x/text => golang.org/x/text v0.14.0
