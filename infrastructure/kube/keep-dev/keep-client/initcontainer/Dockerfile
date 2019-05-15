FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55

ARG CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
ENV CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS=$CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

ARG ETH_HOSTNAME
ENV ETH_HOSTNAME=$ETH_HOSTNAME

ARG ETH_HOST_PORT
ENV ETH_HOST_PORT=$ETH_HOST_PORT

ARG ETH_NETWORK_ID
ENV ETH_NETWORK_ID=$ETH_NETWORK_ID

COPY ./StakingProxy.json /tmp/StakingProxy.json
COPY ./TokenStaking.json /tmp/TokenStaking.json
COPY ./KeepToken.json /tmp/KeepToken.json

COPY ./unlock-and-stake-eth-account.js /tmp/unlock-and-stake-eth-account.js

ENTRYPOINT ["node", "./unlock-and-stake-eth-account.js"]
