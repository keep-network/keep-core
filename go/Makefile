build:
	go build github.com/keep-network/keep-core/go/

proto:
	go generate sh -c "protoc --proto_path=$GOPATH/src:. --gogofaster_out=. *.proto"
