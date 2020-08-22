#!/bin/bash

set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ENV=dev
cd ${DIR}/..
go run ./main.go build -i ${DIR}/../_example/example.yml -i ${DIR}/../_example/params_${ENV}.yml -o ${DIR}/../_example/container/example.go
#go run ./main.go dump-params -i ${DIR}/../_example/example.yml -i ${DIR}/../_example/params_${ENV}.yml
cd ${DIR}
#echo TIDY
#go mod tidy
#echo VENDOR
#go mod vendor
#cat go.mod
go run main.go
#cat ${DIR}/../example/container/example.go
