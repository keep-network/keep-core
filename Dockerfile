FROM golang:1.11.4-alpine3.7 AS runtime

ENV APP_NAME=keep-client \
	BIN_PATH=/usr/local/bin

FROM runtime AS gobuild

ENV GOPATH=/go \
	GOBIN=/go/bin \
	APP_NAME=keep-client \
	APP_DIR=/go/src/github.com/keep-network/keep-core \
	BIN_PATH=/usr/local/bin \
	LD_LIBRARY_PATH=/usr/local/lib/

RUN apk add --update --no-cache \
	g++ \
	protobuf \
	git \
	make \
	nodejs \
	python && \
	rm -rf /var/cache/apk/ && mkdir /var/cache/apk/ && \
	rm -rf /usr/share/man

COPY --from=ethereum/solc:0.5.4 /usr/bin/solc /usr/bin/solc

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

# Configure GitHub token to be able to get private repositories.
ARG GITHUBTOKEN
RUN git config --global url."https://$GITHUBTOKEN:@github.com/".insteadOf "https://github.com/"

RUN go get -u github.com/golang/dep/cmd/dep

COPY ./Gopkg.toml ./Gopkg.lock ./
RUN dep ensure -v --vendor-only
RUN cd vendor/github.com/gogo/protobuf/protoc-gen-gogoslick && go install .
RUN cd vendor/github.com/ethereum/go-ethereum/cmd/abigen && go install .

COPY ./contracts/solidity $APP_DIR/contracts/solidity
RUN cd $APP_DIR/contracts/solidity && npm install

COPY ./pkg/net/gen $APP_DIR/pkg/net/gen
COPY ./pkg/chain/gen $APP_DIR/pkg/chain/gen
COPY ./pkg/beacon/relay/entry/gen $APP_DIR/pkg/beacon/relay/entry/gen
COPY ./pkg/beacon/relay/gjkr/gen $APP_DIR/pkg/beacon/relay/gjkr/gen
COPY ./pkg/beacon/relay/dkg/result/gen $APP_DIR/pkg/beacon/relay/dkg/result/gen
COPY ./pkg/beacon/relay/registry/gen $APP_DIR/pkg/beacon/relay/registry/gen
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
