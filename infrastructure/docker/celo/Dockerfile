FROM golang:1.14-alpine AS gobuild

RUN apk add --update --no-cache \
    g++ \
    linux-headers \
    make \
    git

ENV BIN_PATH=/usr/local/bin

# Clone, compile, and install Celo blockchain v1.2.2
RUN git clone https://github.com/celo-org/celo-blockchain.git
WORKDIR /go/celo-blockchain
RUN git checkout v1.2.2
RUN make geth-musl
RUN cp build/bin/geth $BIN_PATH

FROM alpine:3.10

ENV BIN_PATH=/usr/local/bin
# Environment variables used in start.sh
ENV CELO_DATA_DIR=/mnt/data
ENV CELO_INIT_DIR=/celo-init
ENV CELO_RPC_PORT=8545
ENV CELO_WS_PORT=8546

# Volume where chain data are persisted
VOLUME [$CELO_DATA_DIR]

EXPOSE $CELO_RPC_PORT
EXPOSE $CELO_WS_PORT

COPY --from=gobuild $BIN_PATH/geth $BIN_PATH

# Add celo-init directory. It contains genesis.json file and
# keystore used in start.sh
ADD ./celo-init $CELO_INIT_DIR

ADD ./start.sh .
ENTRYPOINT ["./start.sh"]