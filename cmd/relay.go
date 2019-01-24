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

// RelayCommand contains the definition of the relay command-line subcommand and
// its own subcommands.
var RelayCommand cli.Command

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
		cancel()
		return nil
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("Request errored out: [%v].\n", err)
		} else {
			return fmt.Errorf("Request errored for unknown reason.\n")
		}
	}

	return nil
}

// relayEntry requests an entry with a particular id from the threshold relay
// and prints that entry.
func relayEntry(c *cli.Context) error {
	return fmt.Errorf("relay entry lookups are currently unimplemented")
}

// submitRelayEntrySeed creates a new seed entry for the threshold relay, kicking
// off the group selection process, and prints the newly generated value. By
// default, it uses a random big integer as the value.
func submitRelayEntrySeed(c *cli.Context) error {
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

	// Seed the network with the first n bits of pi
	value := relay.GenesisEntryValue()

	entry := &event.Entry{
		RequestID:     big.NewInt(0),
		Value:         value,
		GroupID:       big.NewInt(0),
		PreviousValue: big.NewInt(0),
		Timestamp:     time.Now().UTC(),
	}

	provider.ThresholdRelay().SubmitRelayEntry(
		entry,
	).OnSuccess(func(data *event.Entry) {
		fmt.Fprintf(
			os.Stdout,
			"Submitted relay entry: [%v].\n",
			data,
		)
		wait <- nil
		return
	}).OnFailure(func(err error) {
		wait <- err
		return
	})

	provider.ThresholdRelay().OnRelayEntryGenerated(func(entry *event.Entry) {
		fmt.Fprintf(
			os.Stderr,
			"Relay entry received with value: [%v].\n",
			entry.Value,
		)
	})

	select {
	case err := <-wait:
		if err != nil {
			return fmt.Errorf(
				"Error in submitting relay entry: [%v].\n",
				err,
			)
		} else {
			cancel()
		}
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("Request errored out: [%v].\n", err)
		} else {
			return fmt.Errorf("Request errored for unknown reason.\n")
		}
	}

	return nil
}
