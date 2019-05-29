FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55
RUN npm install concat-stream@2.0.0
RUN npm install toml@3.0.0
RUN npm install tomlify-j0.4@3.0.0

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
COPY ./KeepRandomBeacon.json /tmp/KeepRandomBeacon.json
COPY ./KeepGroup.json /tmp/KeepGroup.json

COPY ./keep-client-bootstrap-peer-template.toml /tmp/keep-client-bootstrap-peer-template.toml
COPY ./keep-client-standard-peer-template.toml /tmp/keep-client-standard-peer-template.toml

COPY ./create-unlock-stake-account.js /tmp/create-unlock-stake-account.js



ENTRYPOINT ["node", "./create-unlock-stake-account.js"]
