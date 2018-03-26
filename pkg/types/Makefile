.PHONY: install proto-gogo

install: proto-darwin proto-gogo

proto-darwin:
		brew install protobuf
		@echo "Installed protobuf compiler for Darwin"

proto-gogo:
		go get github.com/gogo/protobuf/proto \
					 github.com/gogo/protobuf/jsonpb \
				   github.com/gogo/protobuf/protoc-gen-gogo \
           github.com/gogo/protobuf/gogoproto

		@echo "Installed gogo deps"
