package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/urfave/cli"

	relaychain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
)

// APICommand contains the definition of the api command-line subcommand.
var APICommand cli.Command

const apiDescription = `The api command boots up an HTTP server with a JSON API
	with access to a few basic details regarding the current state of the
	relay.`

const (
	defaultHTTPPort = 8080
)

func init() {
	APICommand = cli.Command{
		Name:        "api",
		Usage:       "Provides access to relay details via an HTTP JSON API.",
		Description: apiDescription,
		Action:      startAPIServer,
		Flags: []cli.Flag{
			&cli.IntFlag{
				// Reuse port flags from start command.
				Name:  portFlag + "," + portShort,
				Value: defaultHTTPPort,
			},
		},
	}
}

func easySmokeTest(groupSize, threshold int) chain.Handle {
	minimumStake := 10000

	chainHandle := local.Connect(
		groupSize,
		threshold,
		big.NewInt(int64(minimumStake)),
	)

	context := context.Background()

	for i := 0; i < groupSize; i++ {
		createNode(context, chainHandle, groupSize, threshold)
	}

	// Give the nodes a sec to get going.
	<-time.NewTimer(time.Second).C

	chainHandle.ThresholdRelay().SubmitRelayEntry(&event.Entry{
		RequestID:     big.NewInt(0),
		Value:         big.NewInt(0),
		GroupPubKey:   big.NewInt(0).Bytes(),
		Seed:          big.NewInt(0),
		PreviousEntry: &big.Int{},
	})

	return chainHandle
}

// startAPIServer starts a new HTTP server and blocks until the process is
// terminated.
func startAPIServer(c *cli.Context) error {
	/*config, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}*/

	fmt.Printf("Connecting to chain provider...\n")
	//chainProvider, err := ethereum.Connect(config.Ethereum)
	chainProvider := easySmokeTest(5, 3)

	//if err != nil {
	//	return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	//}

	port := c.Int(portFlag)
	// Force to default if a bad value is passed.
	if port <= 0 {
		port = defaultHTTPPort
	}

	http.HandleFunc(
		"/api/v0/random-beacon/latest",
		func(rw http.ResponseWriter, r *http.Request) {
			serveLatestEntry(chainProvider.ThresholdRelay(), rw, r)
		})

	fmt.Printf("Starting HTTP server at port %v...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

type beaconEntry struct {
	RequestID *big.Int `json:"requestId"`
	Value     *big.Int `json:"value"`
}

func serveLatestEntry(tr relaychain.Interface, rw http.ResponseWriter, r *http.Request) {
	latestRequestID, err := tr.LatestServedRelayRequestID()
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("Error interacting with chain provider: [%v].", err)))
		return
	}
	latestValue, err := tr.LatestServedRelayValue()
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("Error interacting with chain provider: [%v].", err)))
		return
	}

	entry := beaconEntry{latestRequestID, latestValue}
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("Error serializing JSON: [%v].", err)))
	}

	rw.WriteHeader(200)
	rw.Write(jsonBytes)
}
