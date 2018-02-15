#!/usr/bin/env bash
# Build vendor directory from latest source (assumes vendor directory is checked into github)
# Usage: build-keep-client-docker-img.sh [--force]
# Note:  Only use the --force flag locally when testing.  Your last test before pushing to github
#        should be run w/o the --force flag to ensure your latest updates will work in production.
FORCE=$1
if [ "$(basename $(pwd))" != "keep-core" ] || [ ! -d ./go ]; then
	echo "You should run $(basename $0) from the github.com/keep-net/keep-core directory"
	if [ "$FORCE" != "--force" ]; then exit 2;fi
fi

if [ $(git status | grep "On branch master" | wc | awk '{print $1}') == "0" ]; then
	echo "WARNING: You are not on the master branch!  Are you sure you want to continue? (CTRL+C to abort)"
	read x
fi
if [ $(git status | grep "Your branch is up to date with 'origin/master'." | wc | awk '{print $1}') == "0" ]; then
	echo "Your branch (local file system) is NOT up-to-date with the master branch."
	if [ "$FORCE" != "--force" ]; then exit 2;fi
fi
if [ $(git status | grep "nothing to commit, working tree clean" | wc | awk '{print $1}') == "0" ]; then
	echo "Have you committed all of your changes?"
	if [ "$FORCE" != "--force" ]; then exit 2;fi
fi

# Go to source directory
cd go
# Build vendor from current source
dep ensure
# Back to the keep-core directory that has the Dockerfile.
cd ..

IMG=keep-client
DOCKERFILE=Dockerfile

docker build -t "$IMG" -f "$DOCKERFILE" .
