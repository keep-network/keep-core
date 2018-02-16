package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	golog "github.com/ipfs/go-log"
	gologging "github.com/whyrusleeping/go-logging"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/dfinity/go-dfinity-crypto/rand"
	"github.com/ipfs/go-datastore"
	"github.com/keep-network/keep-core/go/node"
	"github.com/libp2p/go-libp2p-host")

var (
	store = datastore.NewMapDatastore()
	ps    = pstore.NewPeerstore()
)

func ConfigHostHandlers(host host.Host) {
	node.InitTestHandler(host)
	node.InitGetHandler(host)
	node.InitAddHandler(host)
}

func main() {

	// LibP2P code uses golog to log messages. They log with different string IDs (i.e. "swarm").
	// We can control the verbosity level for all loggers with:
	golog.SetAllLoggers(gologging.INFO) // Change to DEBUG for extra info

	// Parse options from the command line
	p2pListenPort := flag.Int("p2pListenPort", 0, "Port that listens for incoming connections")
	p2pEncryption := flag.Bool("p2pEncryption", false, "Enable secure IO")
	idGenerationSeed := flag.Int64("idGenerationSeed", 0, "Set random seed for ID generation")
	flag.Parse()

	if p2pListenPort == nil {
		log.Fatal("Please provide a port to bind on with --port")
	}
	// Make a host that listens on the given multiaddress
	host, err := node.MakeBasicHost(ps, *p2pListenPort, *p2pEncryption, *idGenerationSeed)
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

	//sleep long enough for peerlists to get built.
	time.Sleep(time.Duration(6) * time.Second)
	node.Test(host)
	select {}

	r := rand.NewRand()
	fmt.Printf("%v\n", r)
}
