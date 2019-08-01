#!/bin/bash
set -e

echo "Installing golang requirements..."
for pkg in golang dep; do
  brew list $pkg &> /dev/null || brew install $pkg
done

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

echo "Installing bn requirements..."
for pkg in gmp openssl llvm; do
  brew list $pkg &> /dev/null || brew install $pkg
done

export PATH="/usr/local/opt/llvm/bin:${PATH}"
export LD_LIBRARY_PATH="/usr/local/opt/openssl/lib/:${LD_LIBRARY_PATH}"
export DYLD_LIBRARY_PATH="${DYLD_LIBRARY_PATH}:${LD_LIBRARY_PATH}"

echo "Installing command line developer tools..."
xcode-select --install || true

if ! [ -x "$(command -v protoc-gen-gogoslick)" ]; then
  echo 'WARNING: protoc-gen-gogoslick command is not available'
  echo 'WARNING: please check whether $GOPATH/bin is added to your $PATH'
fi

echo "Installing contracts/solidity npm and requirements..."
brew list npm &>/dev/null || brew install npm
cd ../contracts/solidity && npm install && cd ../../scripts

echo "Installing truffle..."
npm install -g truffle

echo "******************************************************************"
echo "*** Please configure PATH, LD_LIBRARY_PATH, DYLD_LIBRARY_PATH  ***"
echo "*** environment variables in your shell configuration file.    ***"
echo "*** Ex. for bash shell - ~/.bash_profile, for zsh - ~/.zshrc   ***"
echo "******************************************************************"

echo "Ready to rock! See above for any extra environment-related instructions."
