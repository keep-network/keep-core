FROM circleci/golang:1.9.2

WORKDIR /go/src/keep-network/beacon/
# RUN go get -u github.com/golang/dep/cmd/dep
COPY ./go/ ./
# RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
