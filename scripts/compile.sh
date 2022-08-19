#!/bin/bash
# vim: ai:ts=8:sw=8:noet
# Compile this bin
# Intended to be run from local machine or CI
set -eufo pipefail
IFS=$'\t\n'

# Check required commands are in place
command -v go >/dev/null 2>&1 || { echo 'please install go or use image that has it'; exit 1; }

cmd_path="$PWD/cmd/notifier"
bin_path="$PWD/bin/notifier"

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o "$bin_path" "$cmd_path"
