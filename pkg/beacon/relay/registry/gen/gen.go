package gen

//go:generate sh -c "protoc --proto_path=$GOPATH/src:. --gogoslick_out=. */*.proto"
