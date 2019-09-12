#!/bin/sh

echo "Building protobufs files ..."
protoc  --proto_path protos/ \
        --go_out=plugins=grpc:. \
        protos/kademlia/networking/internal_api_service.proto

echo "Building go application ..."
go build main.go