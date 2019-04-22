//go:generate sh -c "rm -f ./pb/*pb.go; protoc --proto_path=$GOPATH/src:. --gogoslick_out=. */*.proto"
package types
