package cmd

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	crand "crypto/rand"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/urfave/cli"
)

// RelayCommand contains the definition of the relay command-line subcommand and
// its own subcommands.
var RelayCommand cli.Command

const relayDescription = `The relay command allows interacting with Keep 
   threshold relay implementation of a random beacon. It is possible
   to request a new entry (equivalent to asking the bacon for a new random
   number. The request subcommand waits for the entry to appear on-chain and 
   then reports its value. It is also possible to seed the relay - this action 
   can be done only once and can not be repeated ever again.`

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
				Name:   "seed",
				Usage:  "Seeds the threshold relay. Can be executed only one time",
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

	// seed is a cryptographically secure pseudo-random number in [0, 2^256)
	seed, err := crand.Int(
		crand.Reader,
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
	)
	if err != nil {
		return fmt.Errorf("could not generate seed: [%v]", err)
	}

	requestMutex := sync.Mutex{}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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

	provider.ThresholdRelay().RequestRelayEntry(seed).
		OnComplete(func(request *event.Request, err error) {
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

// submitRelayEntrySeed creates a new seed entry for the threshold relay, kicking
// off the group selection process, and prints the newly generated value.
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
	defer cancel()

	provider.ThresholdRelay().SubmitRelayEntry(
		relay.GenesisRelayEntry(),
	).OnComplete(func(data *event.Entry, err error) {
		if err != nil {
			wait <- err
			return
		}
		fmt.Printf("Submitted seed relay entry: [%+v]\n", data)
		wait <- nil
		return
	})

	select {
	case err := <-wait:
		if err != nil {
			return fmt.Errorf("error in submitting seed relay entry: [%v]", err)
		}
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			return fmt.Errorf("context done with error: [%v]", err)
		}
		return nil
	}
	return nil
}
