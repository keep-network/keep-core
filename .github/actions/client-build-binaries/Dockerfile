FROM golang:1.18.3-alpine3.16 AS gobuild

RUN apk add --update --no-cache \
    g++ \
    linux-headers \
    protobuf \
    git \
    make \
    nodejs \
    npm \
    bash \
    python3 \
    tar \
    jq && \
    rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
    rm -rf /usr/share/man

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
