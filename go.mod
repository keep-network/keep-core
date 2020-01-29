module github.com/keep-network/keep-core

go 1.12

replace (
	github.com/libp2p/go-libp2p-pubsub => github.com/keep-network/go-libp2p-pubsub v0.0.3-0.20200121091942-499109b16542
	github.com/urfave/cli => github.com/keep-network/cli v1.20.0
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/ethereum/go-ethereum v1.9.7
	github.com/gogo/protobuf v1.3.1
	github.com/ipfs/go-datastore v0.1.1
	github.com/ipfs/go-log v0.0.1
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20190611114437-e92bd71e8199
	github.com/keep-network/keep-common v0.1.1-0.20191203134929-648c427de66e
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-libp2p v0.4.1
	github.com/libp2p/go-libp2p-connmgr v0.1.0
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/libp2p/go-libp2p-kad-dht v0.3.0
	github.com/libp2p/go-libp2p-peerstore v0.1.4
	github.com/libp2p/go-libp2p-pubsub v0.2.2
	github.com/libp2p/go-libp2p-secio v0.2.1
	github.com/libp2p/go-yamux v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/multiformats/go-multiaddr v0.2.0
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/urfave/cli v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/tools v0.0.0-20190925230517-ea99b82c7b93
)
