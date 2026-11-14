#!/bin/sh
# Builds, vets and tests one check, then prints its --describe answer, as CI
# runs it for each check that scripts/changed-checks.sh names.
set -eu
cd "$(dirname "$0")/.."

name="${1:?usage: scripts/test-check.sh <check name>}"
pkg="./cmd/fitness-check-$name"
go build -o out/ "$pkg"
go vet "$pkg"
go test -count=1 "$pkg"
"out/fitness-check-$name" --describe
