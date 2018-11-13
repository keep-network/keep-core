FROM ethereum/client-go:alpine AS build
MAINTAINER "Markus Fix <Markus.Fix@keep.network>

# copy our genesis block definition
COPY genesis.json /root/genesis.json
# initialize chain with our genesis.json parameters
RUN /geth init /root/genesis.json
# copy passphrase file
COPY testnet-account-passphrase.txt /root/passphrase
# create new account for this Keep client
RUN /geth account new  --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 > /root/account1
# create new account for Keep peers
RUN /geth account new  --password /root/passphrase | \
    cut -d "{" -f2 | cut -d "}" -f1 > /root/account2
# save accounts for later retrieval
RUN /geth account list >/root/accounts.txt && cat /root/accounts.txt

# cleanup
RUN \
    rm /root/passphrase

# copy entrypoint script
COPY ./docker-entrypoint.sh /

# provision the geth daemon
EXPOSE 8545
EXPOSE 8546
EXPOSE 30303
ENTRYPOINT ["/docker-entrypoint.sh"]
