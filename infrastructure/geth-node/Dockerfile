FROM ethereum/client-go:v1.9.6
MAINTAINER "Thesis.co"

RUN apk add --no-cache --update git
RUN apk add --no-cache --update python
RUN apk add --no-cache --update build-base
RUN apk add --no-cache --update nodejs npm
RUN apk add --no-cache --update bash
RUN apk add --no-cache --update jq
RUN apk add --no-cache --update curl

# configure log rotation

RUN npm install pm2 -g
RUN pm2 install pm2-logrotate
RUN pm2 set pm2-logrotate:max_size 100M
RUN pm2 set pm2-logrotate:compress true
RUN pm2 set pm2-logrotate:rotateInterval '23 * * *'

# install code to report in at the registry of the bootnode (dashboard)
RUN git clone https://github.com/lispmeister/bootnode-registrar.git /root/lib/bootnode
WORKDIR /root/lib/bootnode
RUN npm install

# install ethStatsApi to report local stats to dashboard
RUN git clone https://github.com/lispmeister/eth-net-intelligence-api.git /root/lib/ethStatsApi
WORKDIR /root/lib/ethStatsApi
RUN npm install

# Change to /root before provisioning our services
WORKDIR /root

RUN mkdir .geth

# copy passphrase file
COPY testnet-account-passphrase.txt passphrase

# Copy keystore
# If you need a copy of the keystore it's in /keep-core/private-testnet/keyfles
ADD keystore .geth/keystore

# Create genesis file
COPY genesis-template.json genesis-template.json
COPY geth-init.sh geth-init.sh
RUN /root/geth-init.sh

# provision our three services (check app.json for details)
COPY app.json app.json
COPY run-geth.sh run-geth.sh

ENTRYPOINT ["pm2", "start", "--no-daemon", "app.json"]
