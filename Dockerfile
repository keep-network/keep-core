FROM golang:1.9.4-alpine3.7

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_REPO_DIR=/go/src/github.com/keep-network/keep-core/go \
	APP_NAME=keep-client

RUN apk add --update --no-cache \
	git && \
	mkdir -p /go/src && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man

COPY ./go $APP_REPO_DIR
WORKDIR $APP_REPO_DIR

RUN go get -u github.com/golang/dep/cmd/dep
COPY ./go/Gopkg.toml ./go/Gopkg.lock ./
RUN dep ensure --vendor-only

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME . && \
	mv $APP_NAME /usr/local/bin/

FROM alpine:3.7

ENV BN_VERSION=d1a44d2f242692601b3e150b59044ab82f265b65

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
	mkdir -p /go/src && \
    git clone https://github.com/dfinity/bn /bn && \
	cd /bn && \
    git reset --hard $BN_VERSION && \
	make install && make && \
	rm -rf /bn && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man && \
	apk del git make clang llvm

COPY --from=0 /usr/local/bin/keep-client /usr/local/bin/

ENV LD_LIBRARY_PATH=/usr/local/lib/

ENTRYPOINT ["keep-client"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
