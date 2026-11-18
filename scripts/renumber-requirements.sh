#!/bin/sh
# Renumbers requirements once by the old=new pairs in scripts/renumber-requirements.txt, in one pass.
# Docs and tests move with git mv, and their references follow, but changelog entries and example numbers stay.
set -eu
cd "$(dirname "$0")/.."

RENUMBER_MAP="${1:-scripts/renumber-requirements.txt}"
export RENUMBER_MAP
# load reads the pairs into %m, and n maps a number, keeping one it lacks.
load='BEGIN { open M, $ENV{RENUMBER_MAP} or die "$ENV{RENUMBER_MAP}: $!"; while (<M>) { $m{$1} = $2 while /\b([0-9]{4})=([0-9]{4})\b/g } }
sub n { exists $m{$_[0]} ? $m{$_[0]} : $_[0] }'

# The docs must carry exactly the old numbers, so a second run stops instead of renumbering again.
have="$(ls docs/requirements | grep -oE "^[0-9]{4}" | sort)"
want="$(grep -oE "\b[0-9]{4}=" "$RENUMBER_MAP" | tr -d = | sort)"
[ "$have" = "$want" ] || { echo "docs/requirements does not hold the old numbers in $RENUMBER_MAP; nothing renamed" >&2; exit 1; }

for path in docs/requirements/[0-9][0-9][0-9][0-9]-*.md requirements/[0-9][0-9][0-9][0-9]*_test.go; do
  to="$(printf '%s\n' "$path" | perl -pe "$load"' s/(?<![0-9])([0-9]{4})(?![0-9])/n($1)/ge')"
  [ "$to" = "$path" ] || git mv "$path" "$to.renumbering"
done
for path in docs/requirements/*.renumbering requirements/*.renumbering; do
  [ -e "$path" ] && git mv "$path" "${path%.renumbering}"
done

git ls-files -z | grep -zv -e '^changelog/' -e '^CHANGELOG.md$' -e '^tests/examples/' \
  -e '^scripts/renumber-requirements' -e '^scripts/migrate-checks-from-engine.sh$' -e '^go.sum$' |
  xargs -0 perl -i -pe "$load"'
  if ($ARGV =~ m{^docs/requirements/[0-9]{4}-.*\.md$}) { s/^(#|-) ([0-9]{4})(?=[ .])/"$1 " . n($2)/e }
  if ($ARGV =~ m{^requirements/.*\.go$}) { s/(?:(?<=[A-Za-z])|(?<=checks-))([0-9]{4})(?![0-9])/n($1)/ge }
  s{(?<=Requirement )([0-9]{4})(?![0-9])}{n($1)}ge;
  s{(?<=docs/requirements/)([0-9]{4})(?=-checks-)}{n($1)}ge;
  if ($ARGV eq "docs/architecture/05-user-journey.md") {
    s{(?<=Req )([0-9]{4}(?:, [0-9]{4})*)}{join ", ", sort map { n($_) } split /, /, $1}ge;
  }
'
