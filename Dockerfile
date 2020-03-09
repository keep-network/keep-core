FROM golang:1.13.6-alpine3.10 AS runtime

ENV APP_NAME=keep-client \
	BIN_PATH=/usr/local/bin

FROM runtime AS gobuild

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-client \
	APP_DIR=/go/src/github.com/keep-network/keep-core \
	TEST_RESULTS_DIR=/mnt/test-results \
	BIN_PATH=/usr/local/bin \
	LD_LIBRARY_PATH=/usr/local/lib/ \
	# GO111MODULE required to support go modules
	GO111MODULE=on

RUN apk add --update --no-cache \
	g++ \
	linux-headers \
	protobuf \
	git \
	make \
	nodejs \
	npm \
	python && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

COPY --from=ethereum/solc:0.5.10 /usr/bin/solc /usr/bin/solc

RUN go get gotest.tools/gotestsum

RUN mkdir -p $APP_DIR $TEST_RESULTS_DIR

WORKDIR $APP_DIR

# Configure GitHub token to be able to get private repositories.
ARG GITHUB_TOKEN
RUN git config --global url."https://$GITHUB_TOKEN:@github.com/".insteadOf "https://github.com/"

# Get dependencies.
COPY go.mod $APP_DIR/
COPY go.sum $APP_DIR/

RUN go mod download

# Install code generators.
RUN cd /go/pkg/mod/github.com/gogo/protobuf@v1.3.1/protoc-gen-gogoslick && go install .
RUN cd /go/pkg/mod/github.com/ethereum/go-ethereum@v1.9.10/cmd/abigen && go install .

COPY ./contracts/solidity $APP_DIR/contracts/solidity
RUN cd $APP_DIR/contracts/solidity && npm install

COPY ./pkg/net/gen $APP_DIR/pkg/net/gen
COPY ./pkg/chain/gen $APP_DIR/pkg/chain/gen
COPY ./pkg/beacon/relay/entry/gen $APP_DIR/pkg/beacon/relay/entry/gen
COPY ./pkg/beacon/relay/gjkr/gen $APP_DIR/pkg/beacon/relay/gjkr/gen
COPY ./pkg/beacon/relay/dkg/result/gen $APP_DIR/pkg/beacon/relay/dkg/result/gen
COPY ./pkg/beacon/relay/registry/gen $APP_DIR/pkg/beacon/relay/registry/gen
# Need this to resolve imports in generated Ethereum commands.
COPY ./config $APP_DIR/config
RUN go generate ./.../gen 

COPY ./ $APP_DIR/
RUN go generate ./pkg/gen

RUN GOOS=linux go build -a -o $APP_NAME ./ && \
	mv $APP_NAME $BIN_PATH

FROM runtime

COPY --from=gobuild $BIN_PATH/$APP_NAME $BIN_PATH

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client", "-config", "/keepclient/config.toml"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
