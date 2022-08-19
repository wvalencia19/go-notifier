#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Intended to be run from local machine or CI
set -eufo pipefail
export SHELLOPTS	# propagate set to children by default
IFS=$'\t\n'

command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }
command -v golangci-lint >/dev/null 2>&1 || { echo 'please install golangci-lint or use image that has it'; exit 1; }

golangci-lint run -E goimports

go test -v -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
