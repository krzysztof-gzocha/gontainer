#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR}/..
go run ./main.go build -i ${DIR}/../example/example.yml -o ${DIR}/../example/container/example.go
go run ${DIR}/../example/main.go
cd ${DIR}
