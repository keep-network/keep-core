FROM ethereum/client-go:alpine

# alpine sdk - curl, gcc, make, g++, git
RUN apk add --no-cache \
  alpine-sdk \
  bash \
  vim \
  python \
  net-tools \
  nodejs \
  nodejs-npm \
  lsof \
  jq

RUN npm install -g pm2

ENTRYPOINT ["/bin/bash"]
