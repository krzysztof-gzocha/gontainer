#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go run ./main.go build -i ${DIR}/../example2/example.yml -o ${DIR}/../example2/container.go
#cat ${DIR}/../example/container/example.go
