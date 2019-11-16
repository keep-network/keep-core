module github.com/keep-network/keep-core

go 1.12

replace (
	github.com/ethereum/go-ethereum => github.com/keep-network/go-ethereum v1.8.27
	github.com/urfave/cli => github.com/keep-network/cli v1.20.0
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
	github.com/cespare/cp v1.1.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/ipfs/go-datastore v0.1.1
	github.com/ipfs/go-log v0.0.1
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20190611114437-e92bd71e8199
	github.com/keep-network/keep-common v0.0.0-20191002130723-787318dfe040
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-libp2p v0.4.1
	github.com/libp2p/go-libp2p-connmgr v0.1.0
	github.com/libp2p/go-libp2p-core v0.2.5
	github.com/libp2p/go-libp2p-host v0.1.0
	github.com/libp2p/go-libp2p-kad-dht v0.3.0
	github.com/libp2p/go-libp2p-net v0.1.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.4
	github.com/libp2p/go-libp2p-pubsub v0.2.1
	github.com/libp2p/go-yamux v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/multiformats/go-multiaddr v0.1.2
	github.com/multiformats/go-multihash v0.0.9 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/urfave/cli v0.0.0-00010101000000-000000000000
	go.opencensus.io v0.22.2 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/tools v0.0.0-20190925230517-ea99b82c7b93
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)
