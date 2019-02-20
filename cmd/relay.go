package cmd

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/urfave/cli"
)

const (
	defaultRequestID int = 0
)

// RelayCommand contains the definition of the relay command-line subcommand and
// its own subcommands.
var RelayCommand cli.Command

const (
	requestIDFlag  = "request-id"
	requestIDShort = "r"
)

const relayDescription = `The relay command allows access to the two functions
   possible in the Keep threshold relay implementation of a random
   beacon: requesting a new entry (equivalent to asking the beacon
   for a new random number) and retrieving an existing entry (using
   the request ID). Each of these is a subcommand (respectively,
   request and entry). The request subcommand waits for the entry
   to appear on-chain and then reports its value.`

func init() {
	RelayCommand = cli.Command{
		Name:        "relay",
		Usage:       `Provides access to the Keep threshold relay.`,
		Description: relayDescription,
		Subcommands: []cli.Command{
			{
				Name:   "request",
				Usage:  "Requests a new entry from the relay.",
				Action: relayRequest,
			},
			{
				Name:   "entry",
				Usage:  "Requests the entry associated with the given request id from the relay.",
				Action: relayEntry,
			},
			{
				Name:   "submit",
				Usage:  "Submits a new seed entry to the relay; only for testing.",
				Action: submitRelayEntrySeed,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  requestIDFlag + "," + requestIDShort,
						Value: defaultRequestID,
					},
				},
			},
		},
	}
}

// relayRequest requests a new entry from the threshold relay and prints the
// request id. By default, it also waits until the associated relay entry is
// generated and prints out the entry.
func relayRequest(c *cli.Context) error {
	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: [%v]", err)
	}

	provider, err := ethereum.Connect(cfg.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	requestMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	wait := make(chan struct{})
	var requestID *big.Int
	provider.ThresholdRelay().OnRelayEntryRequested(func(request *event.Request) {
		fmt.Fprintf(
			os.Stderr,
			"Relay entry request submitted with id [%s].\n",
			request.RequestID.String(),
		)
		requestMutex.Lock()
		requestID = request.RequestID
		requestMutex.Unlock()
	})

	provider.ThresholdRelay().OnRelayEntryGenerated(func(entry *event.Entry) {
		requestMutex.Lock()
		defer requestMutex.Unlock()

		if requestID != nil && requestID.Cmp(entry.RequestID) == 0 {
			fmt.Fprintf(
				os.Stderr,
				"Relay entry received with value: [%v].\n",
				entry.Value,
			)

			wait <- struct{}{}
		}
	})

	provider.ThresholdRelay().RequestRelayEntry(
		big.NewInt(0),
		big.NewInt(0),
	).OnComplete(func(request *event.Request, err error) {
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error in requesting relay entry: [%v].\n",
				err,
			)
			return
		}
		fmt.Fprintf(
			os.Stdout,
			"Relay entry requested: [%v].\n",
			request,
		)
	})

	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("request errored out [%v]", err)
		}
		return fmt.Errorf("request errored for unknown reason")

	}
}

// relayEntry requests an entry with a particular id from the threshold relay
// and prints that entry.
func relayEntry(c *cli.Context) error {
	return fmt.Errorf("relay entry lookups are currently unimplemented")
}

// submitRelayEntrySeed creates a new seed entry for the threshold relay, kicking
// off the group selection process, and prints the newly generated value. By
// default, it uses a request ID equal `0`.
func submitRelayEntrySeed(c *cli.Context) error {
	requestID := c.Int(requestIDFlag)

	entry := &event.Entry{
		RequestID:     big.NewInt(int64(requestID)),
		Value:         relay.GenesisEntryValue(),
		GroupID:       big.NewInt(int64(requestID)),
		PreviousEntry: big.NewInt(0),
		Timestamp:     time.Now().UTC(),
		Seed:          big.NewInt(0),
	}

	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: [%v]", err)
	}

	provider, err := ethereum.Connect(cfg.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	var (
		wait        = make(chan error)
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	provider.ThresholdRelay().SubmitRelayEntry(
		entry,
	).OnComplete(func(data *event.Entry, err error) {
		if err != nil {
			wait <- err
			return
		}
		fmt.Fprintf(
			os.Stdout,
			"Submitted relay entry: [%+v].\n",
			data,
		)
		wait <- nil
		return
	})

	select {
	case err := <-wait:
		if err != nil {
			return fmt.Errorf(
				"error in submitting relay entry: [%v]",
				err,
			)
		}
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("context done, with error: [%v]", err)
		}
		return nil
	}
	return nil
}
