#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go run ${DIR}/../main.go build -i example/example.yml -o example/container/example.go
go run ${DIR}/*.go
