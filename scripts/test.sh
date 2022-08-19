#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Check tests
# Intended to be run from local machine or CI
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

go test -race -coverprofile=.test_coverage.txt ./...
go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $3}'
