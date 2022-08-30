#!/bin/bash
set -eou pipefail

ROOT_DIR="$(realpath "$(dirname $0)/../bundles")"

if ! npx eth-helper --version &>/dev/null; then
    printf "eth-helper could not be found; installing... \n"
    npm install nkuba/eth-helper
fi

STAKING_PROVIDER=${1-}
if [ -z "${STAKING_PROVIDER}" ]; then
    read -p "Provide Staking Provider name: " STAKING_PROVIDER
fi

STAKING_PROVIDER_DIR="$(realpath "$ROOT_DIR/$STAKING_PROVIDER")"

if [ -d "$STAKING_PROVIDER_DIR" ]; then
    echo "Directory for $STAKING_PROVIDER already exists."
    exit 1
fi

if [ -z "${KEYFILE_PASSWORD+x}" ]; then
    read -s -r -p "Provide password for key file encryption: " KEYFILE_PASSWORD
    if [ -z "$KEYFILE_PASSWORD" ]; then
        printf "KEYFILE_PASSWORD not set\n"
        exit 1
    fi
    printf "\n"
fi

CONFIG_DIR="$STAKING_PROVIDER_DIR/config"
SECRETS_DIR="$STAKING_PROVIDER_DIR/secret"

KEY_FILE_PATH="$CONFIG_DIR/provider-eth-account-key-file.json"
KEY_FILE_PASSWORD_PATH="$SECRETS_DIR/provider-eth-account-password"
PRIVATE_KEY_FILE_PATH="$SECRETS_DIR/provider-eth-account-private-key"

mkdir $STAKING_PROVIDER_DIR
mkdir $SECRETS_DIR
mkdir $CONFIG_DIR

cd $STAKING_PROVIDER_DIR

echo -n "$KEYFILE_PASSWORD" >"$KEY_FILE_PASSWORD_PATH"

geth account new \
    --keystore ./ \
    --password "$KEY_FILE_PASSWORD_PATH"

mv UTC-* $KEY_FILE_PATH

npx eth-helper extract-private-key \
    -k "$KEY_FILE_PATH" \
    -p "$KEY_FILE_PASSWORD_PATH" \
    -o "$PRIVATE_KEY_FILE_PATH"

asciidoctor ../bundle-guide.adoc -o index.html --doctype book

tar -zcvf keep-test-bundle-$STAKING_PROVIDER.tgz .

printf "A bundle was saved: keep-test-bundle-$STAKING_PROVIDER.tgz"

printf "\n\e[1;32mDONE!\n\n\e[0m"
