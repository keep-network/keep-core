FROM ethereum/client-go:alpine AS build
MAINTAINER "Markus Fix <Markus.Fix@keep.network>

# copy passphrase file
# TODO: this needs to be a runtime parameter
COPY testnet-account-passphrase.txt /root/passphrase

# copy entrypoint script
COPY ./docker-entrypoint.sh /

# provision the geth daemon
EXPOSE 8545
EXPOSE 8546
EXPOSE 30303
ENTRYPOINT ["/docker-entrypoint.sh"]
