#!/bin/bash

go go test -coverprofile=coverage.out ./cmd/... ./pkg/...
go tool cover -html=coverage.out
