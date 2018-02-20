FROM golang:1.9.4-alpine3.7 AS build

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-client \
	APP_DIR=/go/src/github.com/keep-network/keep-client \
	BIN_PATH=/usr/local/bin \
	LD_PATH=/usr/local/lib \
	BN_VERSION=d1a44d2f242692601b3e150b59044ab82f265b65

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
	git reset --hard $BN_VERSION && \
	make install && make && \
	rm -rf /bn && \
	mkdir -p /go/src && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

RUN go get -u github.com/golang/dep/cmd/dep
COPY ./go/Gopkg.toml ./go/Gopkg.lock ./
RUN dep ensure --vendor-only

COPY . $APP_DIR
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME ./go && \
	mv $APP_NAME $BIN_PATH && \
	rm -rf $APP_DIR

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client", "-config",  "/keepclient/config.toml"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []