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

echo "Installing command line developer tools..."
xcode-select --install || true

if ! [ -x "$(command -v protoc-gen-gogoslick)" ]; then
  echo 'WARNING: protoc-gen-gogoslick command is not available'
  echo 'WARNING: please check whether $GOPATH/bin is added to your $PATH'
  exit 1
fi

echo "Installing contracts/solidity npm and requirements..."
brew list npm &>/dev/null || brew install npm
cd ../contracts/solidity && npm install && cd ../../scripts

#
# This section is to install 'bn', Barreto-Naehrig curve implementation and BLS.
# If you don't have this you get the error:
#
#	# github.com/keep-network/keep-core/vendor/github.com/dfinity/go-dfinity-crypto/bls
#	../../vendor/github.com/dfinity/go-dfinity-crypto/bls/bls.go:6:10: fatal error: 'bls/bls.h' file not found
#	#include <bls/bls.h>
#	         ^~~~~~~~~~~
#	1 error generated.
#
echo "Installing bn and it's dependencies..."
echo "  As a part of installing bn the llvm compiler is needed.  Usually this takes a few minutes to install."
brew list gmp &>/dev/null || brew install gmp 
brew list openssl &>/dev/null || brew install openssl 
brew list llvm &>/dev/null || brew install llvm
mkdir -p ~/go/src/github.com/keep-network
cd ~/go/src/github.com/keep-network
if [ -d bn ] ; then
	cd bn
	git pull
else
	git clone https://github.com/keep-network/bn.git
	cd bn
fi
make
make install

echo "Ready to rock! See above for any extra environment-related instructions."
