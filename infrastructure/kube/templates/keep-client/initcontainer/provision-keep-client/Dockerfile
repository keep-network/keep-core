FROM node:11 AS runtime

WORKDIR /tmp

COPY ./package.json /tmp/package.json
COPY ./package-lock.json /tmp/package-lock.json

RUN npm install

ARG CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
ENV CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS=$CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

ARG ETH_HOSTNAME
ENV ETH_HOSTNAME=$ETH_HOSTNAME

ARG ETH_RPC_PORT
ENV ETH_RPC_PORT=$ETH_RPC_PORT

ARG ETH_WS_PORT
ENV ETH_WS_PORT=$ETH_WS_PORT

ARG ETH_NETWORK_ID
ENV ETH_NETWORK_ID=$ETH_NETWORK_ID

COPY ./TokenStaking.json /tmp/TokenStaking.json
COPY ./KeepToken.json /tmp/KeepToken.json
COPY ./KeepRandomBeaconService.json /tmp/KeepRandomBeaconService.json
COPY ./KeepRandomBeaconOperator.json /tmp/KeepRandomBeaconOperator.json

COPY ./keep-client-config-template.toml /tmp/keep-client-config-template.toml

COPY ./provision-keep-client.js /tmp/provision-keep-client.js



ENTRYPOINT ["node", "./provision-keep-client.js"]
