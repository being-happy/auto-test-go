#!/bin/bash

protoc -I=$GOPATH --go_out=plugins=grpc:$GOPATH/src $GOPATH/src/auto-test-go/contract/api.proto
