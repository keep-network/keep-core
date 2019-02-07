FROM lispmeister/geth-dev
MAINTAINER "Markus Fix <Markus.Fix@keep.network>

# configure log rotation
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

# copy passphrase file
COPY testnet-account-passphrase.txt passphrase

# define default for KEEP_ACCOUNTS
ARG KEEP_ACCOUNTS=5
# create mining account, KEEP_ACCOUNTS client accounts and generate genesis.json
COPY genesis-template.json genesis-template.json
COPY geth-init.sh geth-init.sh
RUN /root/geth-init.sh $KEEP_ACCOUNTS

# provision our three services (check app.json for details)
COPY app.json app.json
COPY run-geth.sh run-geth.sh

ENTRYPOINT ["pm2", "start", "--no-daemon", "app.json"]
