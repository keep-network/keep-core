#!/bin/bash

set -e

echo "Installing common tools..."

sudo apt-get update

sudo apt-get install -y \
  curl \
  wget \
  git \
  unzip \
  jq \
  python \
  build-essential

if ! [ -x "$(command -v curl)" ]; then echo "curl installation failed"; exit 1; fi
if ! [ -x "$(command -v wget)" ]; then echo "wget installation failed"; exit 1; fi
if ! [ -x "$(command -v git)" ]; then echo "git installation failed"; exit 1; fi
if ! [ -x "$(command -v unzip)" ]; then echo "unzip installation failed"; exit 1; fi
if ! [ -x "$(command -v jq)" ]; then echo "jq installation failed"; exit 1; fi
if ! [ -x "$(command -v python)" ]; then echo "python installation failed"; exit 1; fi
if ! [ -x "$(command -v make)" ]; then echo "build-essential installation failed"; exit 1; fi

echo "Common tools have been installed successfully!"

echo "Installing Node.js and NPM..."

if ! [ -x "$(command -v nvm)" ]; then
  curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash

  echo 'export NVM_DIR="$HOME/.nvm"' >> ~/.profile
  echo '[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"' >> ~/.profile
  source ~/.profile
fi

nvm install 14.3.0
nvm install 11.15.0
nvm alias default 11.15.0
nvm use default

if ! [ -x "$(command -v node)" ]; then echo "Node installation failed"; exit 1; fi
if ! [ -x "$(command -v npm)" ]; then echo "NPM installation failed"; exit 1; fi

echo "Node.js and NPM have been installed successfully!"

echo "Installing Go..."

GOLANG_PACKAGE=go1.13.4.linux-amd64.tar.gz

curl -O https://storage.googleapis.com/golang/$GOLANG_PACKAGE

tar -xvf $GOLANG_PACKAGE
sudo chown -R root:root ./go
sudo mv go /usr/local

echo 'GOPATH="$HOME/go"' >> ~/.profile
echo 'PATH="$PATH:/usr/local/go/bin:$GOPATH/bin"' >> ~/.profile
source ~/.profile

if ! [ -x "$(command -v go)" ]; then echo "Go installation failed"; exit 1; fi

echo "Go has been installed successfully!"

echo "Installing go-ethereum..."

GETH_PACKAGE=geth-alltools-linux-amd64-1.9.9-01744997.tar.gz

curl -O https://gethstore.blob.core.windows.net/builds/$GETH_PACKAGE

tar -xvf $GETH_PACKAGE
mkdir ./go-ethereum && tar -xzf $GETH_PACKAGE -C ./go-ethereum --strip-components=1
sudo chown -R root:root ./go-ethereum
sudo mv ./go-ethereum/* /usr/local/bin

if ! [ -x "$(command -v geth)" ]; then echo "go-ethereum installation failed"; exit 1; fi

echo "go-ethereum has been installed successfully!"

echo "Installing Solidity..."

SOLC_VERSION=v0.5.17

wget https://github.com/ethereum/solidity/releases/download/$SOLC_VERSION/solc-static-linux

chmod 755 solc-static-linux
sudo mv solc-static-linux /usr/local/bin
sudo ln -s -f /usr/local/bin/solc-static-linux /usr/local/bin/solc

if ! [ -x "$(command -v solc)" ]; then echo "Solidity installation failed"; exit 1; fi

echo "Solidity has been installed successfully!"

echo "Installing Protobuf..."

PROTOC_VERSION=3.11.4
PROTOC_PACKAGE=protoc-$PROTOC_VERSION-linux-x86_64.zip

wget https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOC_VERSION/$PROTOC_PACKAGE

mkdir ./protoc && unzip $PROTOC_PACKAGE -d ./protoc
chmod 755 -R ./protoc
sudo mv protoc/bin/protoc /usr/local/bin
sudo mv protoc/include/* /usr/local/include

go install github.com/gogo/protobuf/protoc-gen-gogoslick@latest

if ! [ -x "$(command -v protoc)" ]; then echo "protoc installation failed"; exit 1; fi
if ! [ -x "$(command -v protoc-gen-gogoslick)" ]; then echo "protoc-gen-gogoslick installation failed"; exit 1; fi

echo "Protobuf has been installed successfully!"

echo "Installing Truffle..."

npm install -g truffle@5.0.30

if ! [ -x "$(command -v truffle)" ]; then echo "Truffle installation failed"; exit 1; fi

echo "Truffle has been installed successfully!"