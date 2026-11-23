#!/bin/sh
# Runs a command and writes what it prints as one step output, so a workflow step stays one call.
# Usage: scripts/step-output.sh <name> <command> [argument...]; it appends <name>=<output> to $GITHUB_OUTPUT.
set -eu

name="$1"
shift
value="$("$@")"
echo "$name=$value" >> "$GITHUB_OUTPUT"
