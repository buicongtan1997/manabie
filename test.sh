#!/usr/bin/env bash
go test -c pkg/services/v1/product_service_test.go
MANABIE_ENV=test ./v1.test
rm ./v1.test

go test -c test/api_test.go
MANABIE_ENV=test ./test.test
rm ./test.test
