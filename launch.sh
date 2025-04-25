#!/bin/sh

#set -e pipefail

if [ -z $DATABASE ]; then
    echo "An error has occurred (DATABASE is not set)."
    exit 1
fi

PWD=$(pwd)
export GOPATH=$PWD/../gopath
export PATH=$GOPATH/bin:$PATH


go mod download
go install github.com/wlbr/mule

go generate --run "mule.*" ./...

export $(cat .env-$DATABASE | xargs)

go run -race http/main/*.go

