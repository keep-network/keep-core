module github.com/keep-network/keep-core

go 1.12

replace github.com/urfave/cli => github.com/keep-network/cli v1.20.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/aristanetworks/goarista v0.0.0-20200206021550-59c4040ef2d3 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/celo-org/celo-blockchain v0.0.0-20210222234634-f8c8f6744526
	github.com/ethereum/go-ethereum v1.10.1
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa
	github.com/google/uuid v1.1.5
	github.com/ipfs/go-datastore v0.4.4
	github.com/ipfs/go-log v1.0.4
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20200423153828-ed815bc50aec
	github.com/keep-network/keep-common v1.4.1-0.20210319095805-ebf46d0b62db
	github.com/libp2p/go-addr-util v0.0.2
	github.com/libp2p/go-libp2p v0.10.3
	github.com/libp2p/go-libp2p-connmgr v0.2.4
	github.com/libp2p/go-libp2p-core v0.6.1
	github.com/libp2p/go-libp2p-kad-dht v0.8.3
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-libp2p-pubsub v0.3.3
	github.com/libp2p/go-libp2p-routing v0.1.0 // indirect
	github.com/libp2p/go-libp2p-secio v0.2.2
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/urfave/cli v1.22.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)
