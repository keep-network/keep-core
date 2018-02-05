//go:generate sh -c "protoc --proto_path=$GOPATH/src:. --gogofaster_out=. *.proto"
package pb
