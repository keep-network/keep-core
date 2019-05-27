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

const relayDescription = `The relay command allows interacting with Keep's
	threshold relay. The "request" subcommand allows for requesting a new entry
	from the relay, which is equivalent to asking for a new random number. This
	subcommand waits for the entry to appear on-chain and then reports the value.
	The "genesis" subcommand submits initial genesis value to the relay. This
	action can be done only once and can not be repeated ever again.`

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
				Name:   "genesis",
				Usage:  "Submits genesis relay entry. Can be executed only one time.",
				Action: submitGenesisRelayEntry,
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
	// 2^256 - 1 (uint256) is the maximum seed value supported by smart contract
	seed, err := crand.Int(
		crand.Reader,
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
	)
	if err != nil {
		return fmt.Errorf("could not generate seed: [%v]", err)
	}

	requestMutex := sync.Mutex{}
	ctx := context.Background()

	wait := make(chan struct{})
	var requestID *big.Int
	provider.ThresholdRelay().OnRelayEntryRequested(func(request *event.Request) {
		// It is possible two relay requests will be generated close to each other.
		// To make sure we catch request ID from the correct one, we compare seed
		// with the one we sent.
		if requestID == nil && seed.Cmp(request.Seed) == 0 {
			fmt.Printf(
				"Relay entry request submitted with id [%s]: [%v]\n",
				request.RequestID.String(),
				request,
			)

			requestMutex.Lock()
			requestID = request.RequestID
			requestMutex.Unlock()
		}
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

	fmt.Printf("Requesting for a new relay entry at [%v]\n", time.Now().String())

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
		})

	select {
	case <-wait:
		return nil
	case <-ctx.Done():
		if err != nil {
			return fmt.Errorf("request errored out [%v]", err)
		}

		return fmt.Errorf("request errored for unknown reason")
	}
}

// submitGenesisRelayEntry submits genesis entry for the threshold relay,
// kicking off protocol to create the first group.
func submitGenesisRelayEntry(c *cli.Context) error {
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
		fmt.Printf("Submitted genesis relay entry: [%+v]\n", data)
		wait <- nil
		return
	})

	select {
	case err := <-wait:
		if err != nil {
			return fmt.Errorf("error in submitting genesis relay entry: [%v]", err)
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
