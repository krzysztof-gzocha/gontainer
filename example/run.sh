#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

go run ./main.go build -i ${DIR}/../example/example.yml -o ${DIR}/../example/container/example.go
PERSON_SALARY="3500" PERSON_POSITION="Chief Technology Officer " go run ${DIR}/../example/main.go
#cat ${DIR}/../example/container/example.go
