#!/usr/bin/env bash

# https://app.codecov.io/gh/jin06/binlogo
#
set -e
echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor); do
    go test -v -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

bash <(curl -s https://codecov.io/bash)

rm coverage.txt
