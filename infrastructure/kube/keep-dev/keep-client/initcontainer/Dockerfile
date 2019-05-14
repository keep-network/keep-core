FROM node:11 AS runtime

WORKDIR /tmp

RUN npm install web3@1.0.0-beta.55

COPY ./unlock-and-stake-eth-account.js /tmp/unlock-and-stake-eth-account.js

ENTRYPOINT ["node", "./unlock-and-stake-eth-account.js"]
