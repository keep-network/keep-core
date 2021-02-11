FROM golang:1.14-alpine AS gobuild

RUN apk add --update --no-cache \
    g++ \
    linux-headers \
    make \
    git

ENV BIN_PATH=/usr/local/bin \
    DATA_DIR=/mnt/data

VOLUME [$DATA_DIR]

# Clone, compile, and install Celo blockchain v1.2.2
RUN git clone https://github.com/celo-org/celo-blockchain.git
WORKDIR celo-blockchain
RUN git checkout v1.2.2
RUN make geth-musl
RUN cp build/bin/geth $BIN_PATH

FROM alpine:3.10
ENV BIN_PATH=/usr/local/bin

COPY --from=gobuild $BIN_PATH/geth $BIN_PATH

# Add celo-init directory. It contains genesis.json file and keystore used
# in start.sh
ADD ./celo-init /celo-init

ADD ./start.sh .
ENTRYPOINT ["./start.sh"]