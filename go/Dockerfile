FROM golang:1.9.4-alpine3.7 AS build

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-client \
	APP_REPO_DIR=/go/src/github.com/keep-network/keep-client \
	APP_DIR=/keepclient \
	TOML_CFG_FNAME=config.toml \
	BIN_PATH=/usr/local/bin \
	BN_VERSION=d1a44d2f242692601b3e150b59044ab82f265b65

COPY . $APP_REPO_DIR
WORKDIR $APP_REPO_DIR

RUN apk add --update --no-cache \
	bash \
	clang \
	g++ \
	git \
	gmp \
	gmp-dev \
	libgmpxx  \
	libstdc++ \
	llvm \
	make \
	openssl \
	openssl-dev && \
	git clone https://github.com/dfinity/bn /bn && \
	cd /bn && \
	make install && make && \
	rm -rf /bn && \
	mkdir -p /go/src && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man && \
	mkdir $APP_DIR && \
	cd $APP_REPO_DIR && \
	cp ./$TOML_CFG_FNAME $APP_DIR

RUN go get -u github.com/golang/dep/cmd/dep
COPY ./Gopkg.toml ./Gopkg.lock ./
RUN dep ensure --vendor-only

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME && \
	mv $APP_NAME $BIN_PATH && \
	rm -rf $APP_REPO_DIR


FROM alpine:3.7

ENV	APP_DIR=/keepclient \
	BIN_PATH=/usr/local/bin

RUN mkdir $APP_DIR
COPY --from=build $BIN_PATH/keep-client $BIN_PATH
COPY --from=build $APP_DIR/config.toml  $APP_DIR

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client", "-config",  "/keepclient/config.toml"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
