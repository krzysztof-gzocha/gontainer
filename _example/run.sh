#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ENV=dev
cd ${DIR}/..
go run ./main.go build -i ${DIR}/../_example/example.yml -i ${DIR}/../_example/params_${ENV}.yml -o ${DIR}/../_example/container/example.go
cd ${DIR}
echo TIDY
go mod tidy
echo VENDOR
go mod vendor
cat go.mod
PERSON_SALARY="3500" PERSON_POSITION="Chief Technology Officer " go run main.go
#cat ${DIR}/../example/container/example.go
