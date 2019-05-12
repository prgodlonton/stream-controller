#!/usr/bin/env bash

# create response writer mock
mockgen -package mocks \
    -destination testing/mocks/responsewriter.go \
    net/http \
    ResponseWriter

# create store mock
mockgen -package mocks \
    -destination testing/mocks/store.go \
    github.com/pgodlonton/stream-controller/internal \
    Store
