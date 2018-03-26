1. Run `make install` to grab the protobuf compiler and `protoc-gen-gogo`.

   a. If this fails for any reason...congratulations! You have the protobuf compiler installed.
      Instead run `make proto-gogo` to install the `protoc-gen-gogo` toolchain.

2. If you add or update types, you'll need to autogenerate protobuf code:

   In `$GOPATH/src/github.com/keep-network/keep-core/pkg/types/`, run `go generate`. 
   Notice a new file, `*.pb.go`. This is your generated code. Ensure it's right and check it in.
