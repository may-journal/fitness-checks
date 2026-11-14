#!/bin/sh
# Prints, as a JSON array, the checks a change since the base commit given can
# alter, run inside the repo to judge. A go.mod change or no base names all.
set -eu
cd "$(git rev-parse --show-toplevel)"

base="${1:-}"
checks="$(ls -d cmd/fitness-check-* | sed 's|^cmd/fitness-check-||')"
if [ -z "$base" ] || ! git cat-file -e "$base^{commit}" 2>/dev/null; then
  changed=go.mod
else
  changed="$(git diff --name-only "$base" HEAD)"
fi

module="$(go list -m)"
chosen=""
for check in $checks; do
  if printf '%s\n' "$changed" | grep -qxE 'go\.(mod|sum)'; then
    chosen="$chosen $check"
    continue
  fi
  dirs="$(go list -deps -test -f '{{.ImportPath}}' "./cmd/fitness-check-$check" |
    sed -n "s|^$module/||p" | sed 's| .*||; s|\.test$||' | sort -u)"
  for file in $changed; do
    dir="$(dirname "$file")"
    case "$file" in "cmd/fitness-check-$check/"*) ;; *)
      printf '%s\n' "$dirs" | grep -qx "$dir" || continue ;;
    esac
    chosen="$chosen $check"
    break
  done
done

printf '['
sep=""
for check in $chosen; do
  printf '%s"%s"' "$sep" "$check"
  sep=","
done
printf ']\n'
