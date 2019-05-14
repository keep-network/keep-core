FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55

COPY /tmp/contracts/StakingProxy.json /tmp/contracts/StakingProxy.json
COPY /tmp/contracts/TokenStaking.json /tmp/contracts/TokenStaking.json
COPY /tmp/contracts/KeepToken.json /tmp/contracts/KeepToken.json

COPY ./unlock-and-stake-eth-account.js /tmp/unlock-and-stake-eth-account.js

ENTRYPOINT ["node", "./unlock-and-stake-eth-account.js"]
