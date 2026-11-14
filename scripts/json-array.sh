#!/bin/sh
# Prints its arguments, each a plain name, as one JSON array of strings.
# Scripts hand that list to a workflow matrix: scripts/json-array.sh <name>...
set -eu

printf '['
sep=""
for item in "$@"; do
  printf '%s"%s"' "$sep" "$item"
  sep=","
done
printf ']\n'
