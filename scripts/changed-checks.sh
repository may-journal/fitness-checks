#!/bin/sh
# Prints the checks a change since a base can alter as a JSON array, or all of
# them for go.mod or no base. With --each, it prints one per commit after it.
set -eu
here="$(cd "$(dirname "$0")" && pwd)"
cd "$(git rev-parse --show-toplevel)"

each=""
if [ "${1:-}" = --each ]; then
  each=1
  shift
fi
base="${1:-}"
module="$(go list -m)"
checks="$(ls -d cmd/fitness-check-* | sed 's|^cmd/fitness-check-||')"
for check in $checks; do
  printf '%s\n' "$check" | grep -qxE '[a-z0-9-]+' || {
    echo "changed-checks.sh: cmd/fitness-check-$check is not a check name" >&2
    exit 1
  }
done

if [ -n "$each" ] && ! git cat-file -e "$base^{commit}" 2>/dev/null; then
  echo "changed-checks.sh: --each needs a base commit, not '$base'" >&2
  exit 1
fi

# One go list call pairs each check with the folders of the module packages
# it depends on, tests included.
listed="$(go list -test -f '{{.ImportPath}} {{join .Deps " "}}' ./cmd/...)"
pairs="$(printf '%s\n' "$listed" | awk -v m="$module/" '
  {
    self = $1; sub("^" m, "", self); sub(/\.test$/, "", self)
    # A folder inside a check, such as its requirements tests, belongs to it.
    check = self; sub("^cmd/fitness-check-", "", check); sub(/\/.*/, "", check)
    print check, "cmd/fitness-check-" check
    print check, self
    for (i = 2; i <= NF; i++) if (index($i, m) == 1) {
      dir = substr($i, length(m) + 1); sub(/\.test$/, "", dir)
      print check, dir
    }
  }' | sort -u)"

# owner names the package folder that holds or embeds a file. It is the
# nearest folder up from the file with Go files.
owner() {
  dir="$(dirname "$1")"
  while [ "$dir" != . ] && ! ls "$dir"/*.go >/dev/null 2>&1; do
    dir="$(dirname "$dir")"
  done
  printf '%s\n' "$dir"
}

# json prints the checks the files in $1 can alter as a JSON array.
json() {
  if printf '%s\n' "$1" | grep -qxE 'go\.(mod|sum)'; then
    chosen="$checks"
  else
    owners="$(for file in $1; do owner "$file"; done | sort -u)"
    chosen="$(printf '%s\n' "$pairs" | while read -r check dir; do
      if printf '%s\n' "$owners" | grep -qxF "$dir"; then echo "$check"; fi
    done | sort -u)"
  fi
  "$here/json-array.sh" $chosen
}

if [ -n "$each" ]; then
  for commit in $(git rev-list --reverse "$base..HEAD"); do
    files="$(git diff --name-only "$commit^" "$commit")"
    printf '%s ' "$commit"
    json "$files"
  done
elif [ -z "$base" ] || ! git cat-file -e "$base^{commit}" 2>/dev/null; then
  json go.mod
else
  files="$(git diff --name-only "$base" HEAD)"
  json "$files"
fi
