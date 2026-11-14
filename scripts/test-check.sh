#!/bin/sh
# Builds one check, vets and tests it with the module packages it imports, and
# prints its --describe answer. CI runs it on each check changed-checks names.
set -eu
cd "$(dirname "$0")/.."

name="${1:?usage: scripts/test-check.sh <check name>}"
pkg="./cmd/fitness-check-$name"
go build -o out/ "$pkg"
packages="$(go list -deps -f '{{if .Module}}{{if .Module.Main}}{{.ImportPath}}{{end}}{{end}}' "$pkg")"
go vet $packages
go test -count=1 $packages
"out/fitness-check-$name" --describe
