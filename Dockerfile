FROM golang:1.18.3-alpine3.16 AS build-sources

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
	python3 \
	tar \
	jq && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

RUN go install gotest.tools/gotestsum@latest

RUN mkdir -p $APP_DIR $TEST_RESULTS_DIR

WORKDIR $APP_DIR

# Get dependencies.
COPY go.mod go.sum $APP_DIR/
RUN go mod download

# Install code generators.
RUN cd /go/pkg/mod/google.golang.org/protobuf@v1.28.1/cmd/protoc-gen-go/ && go install .

# Copy source code for generation.
COPY ./pkg/beacon/dkg/result/gen $APP_DIR/pkg/beacon/dkg/result/gen
COPY ./pkg/beacon/entry/gen $APP_DIR/pkg/beacon/entry/gen
COPY ./pkg/beacon/gjkr/gen $APP_DIR/pkg/beacon/gjkr/gen
COPY ./pkg/beacon/registry/gen $APP_DIR/pkg/beacon/registry/gen
COPY ./pkg/chain/ethereum/beacon/gen $APP_DIR/pkg/chain/ethereum/beacon/gen
COPY ./pkg/chain/ethereum/common/gen $APP_DIR/pkg/chain/ethereum/common/gen
COPY ./pkg/chain/ethereum/ecdsa/gen $APP_DIR/pkg/chain/ethereum/ecdsa/gen
COPY ./pkg/chain/ethereum/tbtc/gen $APP_DIR/pkg/chain/ethereum/tbtc/gen
COPY ./pkg/chain/ethereum/threshold/gen $APP_DIR/pkg/chain/ethereum/threshold/gen
COPY ./pkg/net/gen $APP_DIR/pkg/net/gen
COPY ./pkg/tbtc/gen $APP_DIR/pkg/tbtc/gen
COPY ./pkg/tecdsa/dkg/gen $APP_DIR/pkg/tecdsa/dkg/gen
COPY ./pkg/tecdsa/signing/gen $APP_DIR/pkg/tecdsa/signing/gen
COPY ./pkg/tecdsa/gen $APP_DIR/pkg/tecdsa/gen

# If ENVIRONMENT is not set it will download NPM packages versions
# published and tagged as `development`.
ARG ENVIRONMENT=development

COPY ./Makefile $APP_DIR/Makefile
RUN make get_artifacts environment=$ENVIRONMENT

# Need this to resolve imports in generated Ethereum commands.
COPY ./config $APP_DIR/config
RUN make generate environment=$ENVIRONMENT

#
# Build Docker Image
#
FROM build-sources AS build-docker

WORKDIR $APP_DIR

COPY ./ $APP_DIR/

# Client Versioning.
ARG VERSION
ARG REVISION

RUN GOOS=linux make build \
	version=$VERSION \
	revision=$REVISION \
	output_dir=$APP_DIR \
	app_name=$APP_NAME

FROM alpine:3.16 as runtime-docker

ENV APP_NAME=keep-client \
	APP_DIR=/go/src/github.com/keep-network/keep-core \
	BIN_PATH=/usr/local/bin

COPY --from=build-docker $APP_DIR/$APP_NAME $BIN_PATH

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []

#
# Build Binaries
#
FROM build-sources AS build-bins

COPY ./ $APP_DIR/

WORKDIR $APP_DIR

ARG APP_NAME=keep-client
# Client Versioning.
ARG VERSION
ARG REVISION

RUN make build-all version=$VERSION revision=$REVISION output_dir=$APP_DIR/bin app_name=$APP_NAME

FROM scratch as output-bins

ENV APP_DIR=/go/src/github.com/keep-network/keep-core

COPY --from=build-bins $APP_DIR/bin .
