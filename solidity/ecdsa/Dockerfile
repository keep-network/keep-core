FROM node:16-alpine

RUN apk add --update --no-cache \
    git \
    bash \
    && \
    rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
    rm -rf /usr/share/man

ENV WORK_DIR=/ecdsa

RUN mkdir -p $WORK_DIR
WORKDIR $WORK_DIR

# Use `https://` instead of unauthenticated `git://` protocol.
RUN git config --global url."https://".insteadOf git://

COPY package*.json yarn.lock ./
RUN yarn install

COPY . ./

ENTRYPOINT ["npx", "hardhat"]
