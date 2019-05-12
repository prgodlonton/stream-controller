#!/usr/bin/env bash

# run the following to install mockgen:
#   go get github.com/golang/mock/gomock
#   go install github.com/golang/mock/mockgen

# create response writer mock
mockgen -package mocks \
    -destination testing/mocks/response_writer.go \
    net/http \
    ResponseWriter

# create store mock
mockgen -package mocks \
    -destination testing/mocks/store.go \
    github.com/pgodlonton/stream-controller/internal \
    Store
