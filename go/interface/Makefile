
all: genDoc
	go build
	( cd ./CallContract ; make )
	( cd ./WatchForEvent ; make )

genDoc:
	godoc -html . > RelayContract.html

