#!/bin/bash
set -e

torm_version=$(go list -m -f '{{.Version}}' github.com/TechXTT/TORM/v2)
go run "github.com/TechXTT/TORM/v2/cmd/torm@${torm_version}" generate --out-models torm/models
CGO_ENABLED=0 GOOS=linux go build -o main .
