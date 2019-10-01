#!/bin/sh

echo "Compiling protobufs files ..."
protoc --proto_path protos/ \
  --go_out=plugins=grpc:. \
  protos/kademlia/internal_api_service.proto

echo "Compiling golang application ..."
go build -o main.run main.go
