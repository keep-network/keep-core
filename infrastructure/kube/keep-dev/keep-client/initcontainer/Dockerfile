FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55

ARG CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS
ENV CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS=$CONTRACT_OWNER_ETH_ACCOUNT_ADDRESS

ARG ETH_HOSTNAME
ENV ETH_HOSTNAME=$ETH_HOSTNAME

ARG ETH_HOST_PORT
ENV ETH_HOST_PORT=$ETH_HOST_PORT

COPY /tmp/keep-client/contracts/StakingProxy.json /tmp/keep-client/contracts/StakingProxy.json
COPY /tmp/keep-client/contracts/TokenStaking.json /tmp/keep-client/contracts/TokenStaking.json
COPY /tmp/keep-client/contracts/KeepToken.json /tmp/keep-client/contracts/KeepToken.json

COPY ./unlock-and-stake-eth-account.js /tmp/unlock-and-stake-eth-account.js

ENTRYPOINT ["node", "./unlock-and-stake-eth-account.js"]
