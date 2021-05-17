#!/bin/sh

function set_config_string() {
    sed -i "s/\(${1//\//\\/} *= *\).*/\1\"${2//\//\\/}\"/" $3
}

function set_config_object() {
    sed -i "s/\(${1//\//\\/} *= *\).*/\1${2//\//\\/}/" $3
}

if [[ -z "${OPERATOR_ADDRESS}" ]]; then
    echo "OPERATOR_ADDRESS env not set."
    exit 1
fi

if [[ -z "${OPERATOR_KEYFILE}" ]]; then
    echo "OPERATOR_KEYFILE env not set."
    exit 1
fi

if [[ -z "${OPERATOR_DATA_DIR}" ]]; then
    echo "OPERATOR_DATA_DIR env not set."
    exit 1
fi

if [[ -z "${OPERATOR_WALLET_PASSWORD}" ]]; then
    echo "OPERATOR_WALLET_PASSWORD env not set."
    exit 1
fi

if [[ -z "${P2P_PORT}" ]]; then
    echo "P2P_PORT env not set."
    exit 1
fi

if [[ -z "${P2P_PEERS_ARRAY}" ]]; then
    echo "P2P_PEERS_ARRAY env not set."
    exit 1
fi

CORE_CONFIG_FILE="./config.toml"
set_config_string "Address" $OPERATOR_ADDRESS $CORE_CONFIG_FILE
set_config_string "KeyFile" $OPERATOR_KEYFILE $CORE_CONFIG_FILE
set_config_string "DataDir" $OPERATOR_DATA_DIR $CORE_CONFIG_FILE
set_config_object "Port" $P2P_PORT $CORE_CONFIG_FILE
set_config_object "Peers" $P2P_PEERS_ARRAY $CORE_CONFIG_FILE

cd proxy
TARGET_URL=$RSK_NODE_URL \
TARGET_PORT=$RSK_NODE_PORT \
PROXY_PORT=5050 \
HTTP_MODE=0 \
MATCH_REQUESTS=0 \
MUTE_LOGGING=0 \
pm2 start eth.js --name eth-rsk-proxy
cd ..

KEEP_ETHEREUM_PASSWORD=$OPERATOR_WALLET_PASSWORD LOG_LEVEL="debug" keep-app --config $CORE_CONFIG_FILE start