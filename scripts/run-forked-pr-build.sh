#!/usr/bin/env bash
if [ $# -eq 0 ]; then
  echo "You must specify at least one PR number as an argument."
  echo ""
  echo "Usage:"
  echo "	scripts/run-forked-pr-build.sh <PR number>+"
  echo ""
  echo "	<PR number>: The number of a PR on the origin upstream repository."
  exit 1
fi


for PR_NUMBER in "$@"; do
  echo "Fetching and pushing $PR_NUMBER"
  git fetch origin refs/pull/$PR_NUMBER/head:refs/forked-prs/$PR_NUMBER
  git push origin refs/forked-prs/$PR_NUMBER:refs/heads/forked-prs/$PR_NUMBER
done
