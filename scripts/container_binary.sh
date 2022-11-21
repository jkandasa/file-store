#!/bin/sh

source ./scripts/version.sh

# build server binary, exclude web. it is ok to keep web assets separately inside container
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o file-store-server -ldflags "$LD_FLAGS" cmd/server/main.go
