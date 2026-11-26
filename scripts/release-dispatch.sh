#!/bin/sh
# Tells the org's .github repo that a check release published, so it reruns the fitness suite on open pull requests.
# Run it with GH_TOKEN set to a token that may dispatch to that repo; DISPATCH_REPOSITORY names another target.
set -eu
repo="${DISPATCH_REPOSITORY:-may-journal/.github}"
gh api -X POST "repos/$repo/dispatches" -f event_type=fitness-check-released --silent
