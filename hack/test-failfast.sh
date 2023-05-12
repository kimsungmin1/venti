#!/bin/bash
set -euo pipefail
cd $(dirname $0)/../

echo + go mod tidy -v
       go mod tidy -v

cp docs/examples/datasources.test.yml etc/datasources.yaml

for package in $(go list ./...); do
    echo ====================== $package ======================
    echo + go test -race -failfast -v $package 
           go test -race -failfast -v $package
done
