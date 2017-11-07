Generate proto
protoc -I=. --go_out=plugins=grpc:. fibonacci.proto