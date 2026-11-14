#!/bin/sh
# Prints, as a JSON array, the checks a change since the base commit given can
# alter, run inside the repo to judge. A go.mod change or no base names all.
set -eu
cd "$(git rev-parse --show-toplevel)"

base="${1:-}"
module="$(go list -m)"
checks="$(ls -d cmd/fitness-check-* | sed 's|^cmd/fitness-check-||')"
for check in $checks; do
  printf '%s\n' "$check" | grep -qxE '[a-z0-9-]+' || {
    echo "changed-checks.sh: cmd/fitness-check-$check is not a check name" >&2
    exit 1
  }
done

if [ -z "$base" ] || ! git cat-file -e "$base^{commit}" 2>/dev/null; then
  changed=go.mod
else
  changed="$(git diff --name-only "$base" HEAD)"
fi

# One go list call pairs each check with the folders of the module packages
# it depends on, tests included.
listed="$(go list -test -f '{{.ImportPath}} {{join .Deps " "}}' ./cmd/...)"
pairs="$(printf '%s\n' "$listed" | awk -v m="$module/" '
  {
    check = $1; sub("^" m "cmd/fitness-check-", "", check); sub(/\.test$/, "", check)
    print check, "cmd/fitness-check-" check
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

if printf '%s\n' "$changed" | grep -qxE 'go\.(mod|sum)'; then
  chosen="$checks"
else
  owners="$(for file in $changed; do owner "$file"; done | sort -u)"
  chosen="$(printf '%s\n' "$pairs" | while read -r check dir; do
    if printf '%s\n' "$owners" | grep -qxF "$dir"; then echo "$check"; fi
  done | sort -u)"
fi

printf '['
sep=""
for check in $chosen; do
  printf '%s"%s"' "$sep" "$check"
  sep=","
done
printf ']\n'
