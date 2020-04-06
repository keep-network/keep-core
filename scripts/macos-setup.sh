#!/bin/bash
set -e

echo "Installing golang requirements..."
brew list golang &> /dev/null || brew install golang

echo "Installing ethereum requirements..."
brew tap ethereum/ethereum
brew list geth &>/dev/null || brew install geth
brew list solidity &>/dev/null || brew install solidity

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

echo "Installing solidity npm and requirements..."
brew list npm &>/dev/null || brew install npm
cd ../solidity && npm install && cd ../scripts

echo "Installing truffle..."
npm install -g truffle

if ! [ -x "$(command -v protoc-gen-gogoslick)" ]; then
  echo 'WARNING: protoc-gen-gogoslick command is not available'
  echo 'WARNING: please check whether $GOPATH/bin is added to your $PATH'
fi

echo "Ready to rock! See above for any extra environment-related instructions."
