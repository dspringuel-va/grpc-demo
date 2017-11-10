rm ./protos/fibonacci.pb.go
protoc -I=. --go_out=plugins=grpc:. ./protos/fibonacci.proto