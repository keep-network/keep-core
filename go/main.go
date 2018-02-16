package main

import (
	"fmt"
	"log"
	"time"
	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/dfinity/go-dfinity-crypto/rand"
	"github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p-host"
	"github.com/keep-network/keep-core/node"
	. "github.com/keep-network/keep-core/libs/config"
)

var (
	store = datastore.NewMapDatastore()
	ps    = pstore.NewPeerstore()
)

func ConfigHostHandlers(host host.Host) {
	node.InitTestHandler(host)
	node.InitGetHandler(host)
	node.InitAddHandler(host)
}

func init() {
	GetOptions()
}

func main() {

	// LibP2P code uses golog to log messages. They log with different string IDs (i.e. "swarm").
	// We can control the verbosity level for all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Make a host that listens on the given multiaddress
	host, err := node.MakeBasicHost(ps, Config.P2pListenPort, Config.EnableP2pEncryption, Config.IdGenerationSeed)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("host.ID: %s", host.ID()))
	ConfigHostHandlers(host)

	node.BootstrapConnect(host)

	go func() {
		for {
			log.Printf("Host IP: %s, ID: %s, Peers: %v\n", node.GetOutboundIP(), host.ID(), host.Peerstore().Peers())
			time.Sleep(time.Duration(5) * time.Second)
			node.AddPeers(host)
		}
	}()

	// Sleep long enough for peerlists to get built.
	time.Sleep(time.Duration(6) * time.Second)
	node.Test(host)
	select {}

	r := rand.NewRand()
	fmt.Printf("%v\n", r)
}
