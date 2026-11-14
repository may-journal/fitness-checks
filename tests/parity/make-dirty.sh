#!/bin/sh
# Writes the dirty sample repo into the folder named by its argument: files that break many checks on purpose.
# cspell:ignore recieve nothere
# It is a script, not a folder, so this repo's own checks never judge the faults.
set -eu
cd "$1"
mkdir -p docs plans changelog .github/workflows src
cat > '.fitnessrc.json' <<'EOF'
{}
EOF
cat > '.github/workflows/ci.yml' <<'EOF'
name: CI
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: ./.github/actions/missing
      - run: |
          echo one
          echo two
EOF
cat > '.gitignore' <<'EOF'
build/
secret.txt
EOF
cat > 'changelog/bad.md' <<'EOF'
Not an entry
EOF
cat > 'docs/BadName.md' <<'EOF'
# Odd name

text
EOF
cat > 'go.mod' <<'EOF'
module example.test/dirty

go 1.26
EOF
cat > 'main.go' <<'EOF'
package main
import "fmt"
func main(){
x:=1
if x>0&&x<5||x>10{
for i:=0;i<3;i++{
if i==1{ if i==2{ fmt.Println("a") } else if i==3 { fmt.Println("b") } else if i==4 { fmt.Println("c") } }
}
}
// x := compute(1, 2)
// y := compute(3, 4)
}
EOF
cat > 'out.o' <<'EOF'
a
EOF
cat > 'package.json' <<'EOF'
{"name":"dirty","version":"1.0.0","dependencies":{"left-pad":"0.0.1"}}
EOF
cat > 'plans/plan.md' <<'EOF'
# A plan
EOF
cat > 'README.md' <<'EOF'
# Dirty sample

This sentence is not about teh recieve of **bold** text but about something longer that keeps going on and on without any sign of stopping because it has far too many words in a single sentence for a reader to follow.

It is not a plan, it is a test. See [missing](docs/nothere.md).
EOF
cat > 'src/app.js' <<'EOF'
/* eslint-disable */
var a = 1
console.log(a)
EOF
