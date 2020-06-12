module github.com/keep-network/keep-core

go 1.12

replace github.com/urfave/cli => github.com/keep-network/cli v1.20.0

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/VictoriaMetrics/fastcache v1.5.7 // indirect
	github.com/aristanetworks/goarista v0.0.0-20200206021550-59c4040ef2d3 // indirect
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/ethereum/go-ethereum v1.9.10
	github.com/gogo/protobuf v1.3.1
	github.com/google/gofuzz v1.1.0
	github.com/ipfs/go-datastore v0.1.1
	github.com/ipfs/go-log v0.0.1
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20200423153828-ed815bc50aec
	github.com/keep-network/keep-common v1.1.1-0.20200612111801-12829c0d1e0f
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-libp2p v0.4.1
	github.com/libp2p/go-libp2p-connmgr v0.1.0
	github.com/libp2p/go-libp2p-core v0.3.0
	github.com/libp2p/go-libp2p-kad-dht v0.3.0
	github.com/libp2p/go-libp2p-peerstore v0.1.4
	github.com/libp2p/go-libp2p-pubsub v0.2.6-0.20200127182502-25c434f5f772
	github.com/libp2p/go-libp2p-secio v0.2.1
	github.com/libp2p/go-yamux v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/multiformats/go-multiaddr v0.2.0
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/urfave/cli v1.22.1
	golang.org/x/crypto v0.0.0-20200208060501-ecb85df21340
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
)
