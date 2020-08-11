#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

ENV=dev
go run ./main.go build -i ${DIR}/../_example/example.yml -i ${DIR}/../_example/params_${ENV}.yml -o ${DIR}/../_example/container/example.go
PERSON_SALARY="3500" PERSON_POSITION="Chief Technology Officer " go run ${DIR}/../_example/main.go
#cat ${DIR}/../example/container/example.go
