#!/bin/bash
set -e

echo "Installing golang requirements..."
for pkg in golang dep; do
  brew list $pkg &> /dev/null || brew install $pkg
done

echo "Installing ethereum requirements..."
brew list geth &>/dev/null || brew install geth

echo "Installing protobuf requirements..."
# Protobuf
brew list protobuf &>/dev/null || brew install protobuf
go get -u github.com/gogo/protobuf/protoc-gen-gogoslick

echo "Installing precommit requirements..."
brew list pre-commit &>/dev/null || brew install pre-commit
go get -u golang.org/x/tools/cmd/goimports
go get -u golang.org/x/lint/golint

echo "Installing pre-commit and specified hooks..."
pre-commit install --install-hooks

echo "Installing bn requirements..."
for pkg in gmp openssl llvm; do
  brew list $pkg &> /dev/null || brew install $pkg
done

echo "Installing command line developer tools..."
xcode-select --install || true

echo "Installing build versioning utility..."
go get -u https://github.com/stretchr/version

echo "Ready to rock! See above for any extra environment-related instructions."
