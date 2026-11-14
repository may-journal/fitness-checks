#!/bin/sh
# Writes each check's README from its doc.go. With -check it writes nothing
# and fails naming each README that differs, as CI runs it.
set -eu
cd "$(dirname "$0")/.."
exec go run ./tools/readmes "$@"
