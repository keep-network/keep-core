FROM alpine

ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin
ENV REPO="$GOPATH/src/github.com/l3x/relay"

RUN apk add --update --no-cache bash \
            git \
            make \
            clang \
            g++ \
            go && \
            mkdir -p $REPO && \
            mkdir -p $GOPATH/src && \
            rm -rf /var/cache/apk && \
            rm -rf /usr/share/man && \
            apk del git clang

COPY . $REPO
WORKDIR $REPO
RUN go install && rm -rf $GOPATH/src/

ENTRYPOINT ["relay"]
