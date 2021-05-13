FROM golang:1.13.6-alpine3.10 AS gobuild

ARG VERSION
ARG REVISION

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-app \
	APP_DIR=/go/src/keep-core \
	BIN_PATH=/usr/local/bin \
	LD_LIBRARY_PATH=/usr/local/lib/ \
	GO111MODULE=on

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
	protobuf \
	git \
	make \
	python && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

COPY --from=ethereum/solc:0.5.17 /usr/bin/solc /usr/bin/solc

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR
COPY . $APP_DIR/

RUN GOOS=linux go build -ldflags "-X main.version=$VERSION -X main.revision=$REVISION" -a -o $APP_NAME ./ && \
	mv $APP_NAME $BIN_PATH

FROM node:15-alpine AS app

ENV APP_NAME=keep-app \
	BIN_PATH=/usr/local/bin

RUN apk add --update --no-cache git
	# git \
	# nodejs \
	# npm

RUN npm i -g pm2

COPY --from=gobuild $BIN_PATH/$APP_NAME $BIN_PATH

COPY ./configs/config.local.1.toml ./config.toml
COPY entrypoint.sh .

RUN git clone https://github.com/rumblefishdev/tbtc-rsk-proxy.git proxy
RUN cd proxy/node-http-proxy && npm install
RUN cd proxy && npm install

RUN mkdir /data

ENTRYPOINT ["./entrypoint.sh"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
