api/api.pb.go: proto/api.proto
	protoc -I proto/ proto/api.proto --go_out=plugins=grpc:api

all: api/api.pb.go