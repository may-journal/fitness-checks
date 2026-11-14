#!/bin/sh
# Migrates the checks from the fitness-runner engine into this checks monorepo
# as a plain copy of one engine commit, run once per migration.
set -eu

clone="${1:?usage: scripts/migrate-checks-from-engine.sh <fitness-runner clone> <commit>}"
commit="${2:?usage: scripts/migrate-checks-from-engine.sh <fitness-runner clone> <commit>}"
here="$(cd "$(dirname "$0")/.." && pwd)"
work="$(mktemp -d)"
trap 'rm -rf "$work"' EXIT

# local-branch judges only with the engine's branch metadata, so it stays in
# fitness-runner for now, with its two requirement docs. Doc 0150 proves a
# check of fitness-runner's own, in its fitness folder, so it stays too.
deferred_checks="local-branch"
deferred_docs="0067 0068 0150"

# The check-only helpers: no engine package imports them.
helpers="architecture clonedetect mdfilename requirement sharedconf spell
lang/golang/coverage lang/java/checkstyle lang/java/spotless lang/java/version
lang/node/biome lang/node/vitestconf lang/swift/localized lang/swift/swiftlint"

git -C "$clone" archive "$commit" go/cmd go/internal go/requirements docs/requirements tests/examples |
  tar -x -C "$work"

# copy puts the engine's path $1 at this repo's path $2.
copy() {
  mkdir -p "$here/$(dirname "$2")"
  rm -rf "${here:?}/$2"
  cp -R "$work/$1" "$here/$2"
}

for dir in "$work"/go/cmd/fitness-check-*; do
  name="${dir##*/fitness-check-}"
  case " $deferred_checks " in *" $name "*) continue ;; esac
  copy "go/cmd/fitness-check-$name" "cmd/fitness-check-$name"
done
for helper in $helpers; do
  copy "go/internal/$helper" "internal/$helper"
done
for doc in "$work"/docs/requirements/[0-9][0-9][0-9][0-9]-checks-*.md; do
  id="$(basename "$doc" | cut -c1-4)"
  case " $deferred_docs " in *" $id "*) continue ;; esac
  copy "docs/requirements/${doc##*/}" "docs/requirements/${doc##*/}"
  for test in $(cd "$work" && grep -l "^func Test${id}_" go/requirements/*_test.go); do
    copy "$test" "requirements/${test##*/}"
  done
done
for test in "$work"/go/requirements/*_test.go; do
  case "${test##*/}" in [0-9]*) continue ;; esac
  copy "go/requirements/${test##*/}" "requirements/${test##*/}"
done
copy tests/examples tests/examples

cd "$here"
old=github.com/may-journal/fitness-runner/go
new=github.com/may-journal/fitness-checks

# Each check-only helper now lives in this module's internal folder.
for pkg in $(cd internal && find . -name '*.go' ! -name '*_test.go' -exec dirname {} \; | sed 's|^\./||' | sort -u); do
  grep -rl --include='*.go' "\"$old/internal/$pkg\"" . |
    xargs -r perl -pi -e "s|\"\Q$old/internal/$pkg\E\"|\"$new/internal/$pkg\"|g"
done

# Each check becomes package main. Its main hands internal/program the doc.go
# source, so --describe prints its summary and protocol.
for dir in cmd/fitness-check-*; do
  perl -pi -e 's/^package check$/package main/' "$dir"/*.go
  perl -pi -e 's/^func Main\(\) \{$/func main() {/; s/\bcheckkit\.Main\(/program.Main(docSource, /' "$dir"/*.go
  grep -q 'program"' "$dir/main.go" ||
    perl -pi -e "s|^(\t\"$old/pkg/checkkit\")\$|\t\"$new/internal/program\"\n\$1|" "$dir/main.go"
  grep -q '^var docSource string$' "$dir/doc.go" ||
    printf '\nimport _ "embed"\n\n// docSource is this file, whose comment --describe summarizes.\n//\n//go:embed doc.go\nvar docSource string\n' >>"$dir/doc.go"
done
gofmt -w cmd

# Docs name the moved folders by their new paths, and link them here.
paths='s#(?<![\w/.-])go/cmd/fitness-check-#cmd/fitness-check-#g; s#(?<![\w/.-])go/internal/#internal/#g'
links='s#https://github\.com/may-journal/fitness-runner/blob/main/go/cmd/fitness-check-(?!local-branch)#https://github.com/may-journal/fitness-checks/blob/main/cmd/fitness-check-#g; s#https://github\.com/may-journal/fitness-runner/blob/main/go/internal/(architecture|clonedetect|mdfilename|requirement|sharedconf|spell)#https://github.com/may-journal/fitness-checks/blob/main/internal/$1#g; s#https://github\.com/may-journal/fitness-runner/blob/main/docs/requirements/(\d{4})(?=-checks-)#https://github.com/may-journal/fitness-checks/blob/main/docs/requirements/$1#g'
grep -rlE --include='*.md' 'go/(internal|cmd)/|fitness-runner/blob' cmd internal docs | xargs -r perl -pi -e "$links; $paths"
# In Go only comments change, so every check prints what the bundled one did.
grep -rlE --include='*.go' 'go/(internal|cmd)/|fitness-runner/blob' cmd internal requirements |
  xargs -r perl -pi -e "if (m{^\s*//}) { $links; $paths }"

# A README link that no longer resolves points where it pointed from
# go/cmd/<check>: at the moved file here, else at fitness-runner.
for readme in cmd/fitness-check-*/README.md; do
  README="$readme" perl -pi -e '
    use File::Basename;
    sub tidy { my @out; for (split m{/}, shift) { $_ eq ".." ? pop @out : push @out, $_ } join "/", @out }
    sub up { join "/", (".." ) x (() = $_[0] =~ m{[^/]+}g) }
    my $dir = dirname($ENV{README});
    s{\]\((\.\./[^)#\s]+)}{
      my $link = $1;
      my $old = tidy("go/$dir/$link");
      (my $new = $old) =~ s{^go/(cmd|internal)/}{$1/};
      -e "$dir/$link" ? "]($link"
        : -e $new ? "](" . up($dir) . "/$new"
        : "](https://github.com/may-journal/fitness-runner/blob/main/$old";
    }ge' "$readme"
done
