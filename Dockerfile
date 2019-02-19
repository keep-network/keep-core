FROM golang:1.11.4-alpine3.7 AS runtime

ENV APP_NAME=keep-client \
	BIN_PATH=/usr/local/bin \
	LIB_DIR=/usr/local/lib/ \
	INCLUDE_DIR=/usr/local/include/

RUN apk add --update --no-cache \
	gmp \
	gmp-dev \
	libgmpxx  \
	libstdc++ \
	openssl \
	openssl-dev && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man

FROM runtime AS cbuild

ENV BN_VERSION=1c96f7053ea1ebcdbe9f59ce46c79023ef6f8ba0

RUN apk add --update --no-cache \
	clang \
	g++ \
	git \
	llvm \
	make && \
	rm -rf /var/cache/apk && mkdir /var/cache/apk && \
	rm -rf /usr/share/man

RUN git clone https://github.com/keep-network/bn /bn && \
	cd /bn && \
	git reset --hard $BN_VERSION && \
	make install && make && \
	rm -rf /bn

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

COPY --from=cbuild $LIB_DIR $LIB_DIR
COPY --from=cbuild $INCLUDE_DIR $INCLUDE_DIR
COPY --from=ethereum/solc:0.4.24 /usr/bin/solc /usr/bin/solc

RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

RUN go get -u github.com/golang/dep/cmd/dep

COPY ./Gopkg.toml ./Gopkg.lock ./
RUN dep ensure -v --vendor-only
RUN cd vendor/github.com/gogo/protobuf/protoc-gen-gogoslick && go install .
RUN cd vendor/github.com/ethereum/go-ethereum/cmd/abigen && go install .


COPY ./contracts/solidity $APP_DIR/contracts/solidity
RUN cd $APP_DIR/contracts/solidity && npm install

COPY ./pkg/net/gen $APP_DIR/pkg/net/gen
COPY ./pkg/chain/gen $APP_DIR/pkg/chain/gen
COPY ./pkg/beacon/relay/thresholdsignature/gen $APP_DIR/pkg/beacon/relay/thresholdsignature/gen
COPY ./pkg/beacon/relay/gjkr/gen $APP_DIR/pkg/beacon/relay/gjkr/gen
RUN go generate ./.../gen 

COPY ./ $APP_DIR/
RUN go generate ./pkg/gen

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $APP_NAME ./ && \
	mv $APP_NAME $BIN_PATH

FROM runtime

COPY --from=gobuild $BIN_PATH/$APP_NAME $BIN_PATH
COPY --from=cbuild $LIB_DIR $LIB_DIR
COPY --from=cbuild $INCLUDE_DIR $INCLUDE_DIR

# ENTRYPOINT cant handle ENV variables.
ENTRYPOINT ["keep-client", "-config", "/keepclient/config.toml"]

# docker caches more when using CMD [] resulting in a faster build.
CMD []
