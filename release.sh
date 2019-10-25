#!/bin/bash
TAG=$1
BUCKET=$2
if [ "${TAG}" == "" ]; then
    echo Usage: $0 TAG — for example \'v0.1.0\'
    exit 1
fi
go test ./...
OS=linux
ARCH=amd64
BINARY=uptimerobot-${TAG}-${OS}-${ARCH}
GOOS=$OS GOARCH=$ARCH go build -o ./$BINARY
s3cmd put $BINARY $BUCKET
s3cmd setacl $BUCKET/$BINARY --acl-public
