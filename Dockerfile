FROM golang:1.9.4-alpine3.7

ENV GOPATH=/go \
	GOBIN=/go/bin \
    APP_REPO_DIR=/go/src/keep-network/keep-client \
	APP_NAME=keep-client

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
	apk del git make clang llvm && \
	mkdir -p $APP_REPO_DIR

ENV BN_VERSION=d1a44d2f242692601b3e150b59044ab82f265b65

COPY ./go $APP_REPO_DIR
WORKDIR $APP_REPO_DIR

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME . && \
	mv $APP_NAME /usr/local/bin && \
   	rm -rf $APP_REPO_DIR

ENTRYPOINT ["keep-client"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
