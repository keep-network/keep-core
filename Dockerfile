FROM golang:1.18.3-alpine3.16 AS gobuild

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-client \
	APP_DIR=/go/src/github.com/keep-network/keep-core \
	TEST_RESULTS_DIR=/mnt/test-results \
	BIN_PATH=/usr/local/bin \
	LD_LIBRARY_PATH=/usr/local/lib/

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
	protobuf \
	git \
	make \
	nodejs \
	npm \
	bash \
	python3 && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

RUN go install gotest.tools/gotestsum@latest

RUN mkdir -p $APP_DIR $TEST_RESULTS_DIR

WORKDIR $APP_DIR

# Get dependencies.
COPY go.mod go.sum $APP_DIR/
RUN go mod download

# Install code generators.
RUN cd /go/pkg/mod/github.com/gogo/protobuf@v1.3.2/protoc-gen-gogoslick && go install .

COPY ./pkg/net/gen $APP_DIR/pkg/net/gen
COPY ./pkg/chain/common/gen $APP_DIR/pkg/chain/common/gen
COPY ./pkg/chain/ecdsa/gen $APP_DIR/pkg/chain/ecdsa/gen
COPY ./pkg/chain/random-beacon/gen $APP_DIR/pkg/chain/random-beacon/gen
COPY ./pkg/chain/tbtc-v2/gen $APP_DIR/pkg/chain/tbtc-v2/gen
COPY ./pkg/chain/threshold-network/gen $APP_DIR/pkg/chain/threshold-network/gen
COPY ./pkg/beacon/entry/gen $APP_DIR/pkg/beacon/entry/gen
COPY ./pkg/beacon/gjkr/gen $APP_DIR/pkg/beacon/gjkr/gen
COPY ./pkg/beacon/dkg/result/gen $APP_DIR/pkg/beacon/dkg/result/gen
COPY ./pkg/beacon/registry/gen $APP_DIR/pkg/beacon/registry/gen

# If CONTRACTS_NPM_PACKAGE_TAG is not set it will download NPM packages versions
# published and tagged as `development`.
ARG CONTRACTS_NPM_PACKAGE_TAG

# Need this to resolve imports in generated Ethereum commands.
COPY ./config $APP_DIR/config
RUN go generate ./.../gen

COPY ./ $APP_DIR/
RUN go generate ./pkg/gen

# Client Versioning.
ARG VERSION
ARG REVISION

RUN GOOS=linux go build -ldflags "-X main.version=$VERSION -X main.revision=$REVISION" -a -o $APP_NAME ./ && \
	mv $APP_NAME $BIN_PATH

FROM alpine:3.16

ENV APP_NAME=keep-client \
	BIN_PATH=/usr/local/bin

COPY --from=gobuild $BIN_PATH/$APP_NAME $BIN_PATH

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client", "-config", "/keepclient/config.toml"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
