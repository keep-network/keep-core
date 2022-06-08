module github.com/keep-network/keep-core

go 1.16

replace github.com/urfave/cli => github.com/keep-network/cli v1.20.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/aristanetworks/goarista v0.0.0-20200206021550-59c4040ef2d3 // indirect
	github.com/btcsuite/btcd v0.22.0-beta
	github.com/celo-org/celo-blockchain v0.0.0-20210222234634-f8c8f6744526
	github.com/ethereum/go-ethereum v1.10.15
	github.com/gogo/protobuf v1.3.2
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa
	github.com/google/uuid v1.3.0
	github.com/ipfs/go-datastore v0.4.6
	github.com/ipfs/go-log v1.0.5
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20211001132324-54dddf8aebd4
	github.com/keep-network/keep-common v1.7.1-0.20211012131917-7102d7b9c6a0
	github.com/libp2p/go-addr-util v0.2.0
	github.com/libp2p/go-libp2p v0.20.1
	github.com/libp2p/go-libp2p-connmgr v0.4.0
	github.com/libp2p/go-libp2p-core v0.16.1
	github.com/libp2p/go-libp2p-kad-dht v0.16.0
	github.com/libp2p/go-libp2p-pubsub v0.7.0
	github.com/libp2p/go-libp2p-tls v0.5.0
	github.com/multiformats/go-multiaddr v0.4.0
	github.com/urfave/cli v1.22.1
	golang.org/x/crypto v0.0.0-20210813211128-0a44fdfbc16e
)
