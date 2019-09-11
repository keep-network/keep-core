module github.com/keep-network/keep-core

go 1.12

replace (
	github.com/BurntSushi/toml => github.com/keep-network/toml v0.3.0
	github.com/ethereum/go-ethereum => github.com/keep-network/go-ethereum v1.8.27
	github.com/gogo/protobuf => github.com/keep-network/protobuf v0.0.0-20180801131852-baf2ea5cdb7c
	github.com/urfave/cli => github.com/keep-network/cli v0.0.0-20180226030253-8e01ec4cd3e2
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/allegro/bigcache v0.0.0-20190618191010-69ea0af04088 // indirect
	github.com/aristanetworks/goarista v0.0.0-20190409235741-55bc7be9dd31 // indirect
	github.com/btcsuite/btcd v0.0.0-20190824003749-130ea5bddde3
	github.com/cespare/cp v1.1.1 // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gogo/protobuf v1.2.1
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/gxed/hashland v0.0.1 // indirect
	github.com/hashicorp/go-multierror v1.0.0 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/ipfs/go-cid v0.0.3 // indirect
	github.com/ipfs/go-datastore v0.0.5-0.20190418013242-b19d692f0b56
	github.com/ipfs/go-detect-race v0.0.1 // indirect
	github.com/ipfs/go-log v0.0.1
	github.com/ipfs/go-todocounter v0.0.1 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/keep-network/go-libp2p-bootstrap v0.0.0-20190611114437-e92bd71e8199
	github.com/libp2p/go-addr-util v0.0.1
	github.com/libp2p/go-conn-security v0.1.0
	github.com/libp2p/go-libp2p v0.0.0-2019041893200-f1888d98c45b
	github.com/libp2p/go-libp2p-circuit v0.1.1 // indirect
	github.com/libp2p/go-libp2p-connmgr v0.0.0-20190226230108-808e1bf487d7
	github.com/libp2p/go-libp2p-core v0.2.2 // indirect
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-host v0.1.0
	github.com/libp2p/go-libp2p-interface-connmgr v0.0.0-20190226110100-483442f10797
	github.com/libp2p/go-libp2p-interface-pnet v0.1.0 // indirect
	github.com/libp2p/go-libp2p-kad-dht v0.0.0-20190131020845-7246a3b0f441
	github.com/libp2p/go-libp2p-kbucket v0.2.1 // indirect
	github.com/libp2p/go-libp2p-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-nat v0.0.4 // indirect
	github.com/libp2p/go-libp2p-net v0.1.0
	github.com/libp2p/go-libp2p-netutil v0.0.1 // indirect
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/libp2p/go-libp2p-protocol v0.1.0 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.1.1
	github.com/libp2p/go-libp2p-pubsub-router v0.0.0-20190302015942-41fb0d3d905a // indirect
	github.com/libp2p/go-libp2p-record v0.0.0-20190226223446-0f29df9dd657 // indirect
	github.com/libp2p/go-libp2p-routing v0.1.0 // indirect
	github.com/libp2p/go-libp2p-routing-helpers v0.0.0-20190226233042-cdb43d0f1c87 // indirect
	github.com/libp2p/go-libp2p-swarm v0.2.1 // indirect
	github.com/libp2p/go-libp2p-transport v0.1.0 // indirect
	github.com/libp2p/go-stream-muxer v0.1.0 // indirect
	github.com/libp2p/go-ws-transport v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.9 // indirect
	github.com/miekg/dns v1.1.4 // indirect
	github.com/multiformats/go-multiaddr v0.0.4
	github.com/multiformats/go-multiaddr-dns v0.0.3 // indirect
	github.com/multiformats/go-multihash v0.0.7 // indirect
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.6.0 // indirect
	github.com/satori/go.uuid v1.2.1-0.20180103174451-36e9d2ebbde5 // indirect
	github.com/urfave/cli v0.0.0-00010101000000-000000000000
	github.com/whyrusleeping/base32 v0.0.0-20170828182744-c30ac30633cc // indirect
	github.com/whyrusleeping/go-smux-multiplex v3.0.16+incompatible // indirect
	github.com/whyrusleeping/go-smux-multistream v2.0.2+incompatible // indirect
	github.com/whyrusleeping/go-smux-yamux v2.0.10-0.20190306160942-c236ac3b526b+incompatible // indirect
	github.com/whyrusleeping/mdns v0.0.0-20180901202407-ef14215e6b30 // indirect
	github.com/whyrusleeping/yamux v1.2.0 // indirect
	go.opencensus.io v0.22.1 // indirect
	golang.org/x/crypto v0.0.0-20190911031432-227b76d455e7
	golang.org/x/net v0.0.0-20190909003024-a7b16738d86b // indirect
	golang.org/x/sys v0.0.0-20190910064555-bbd175535a8b // indirect
	golang.org/x/tools v0.0.0-20190418235243-4796d4bd3df0
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)
