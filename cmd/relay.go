package cmd

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	crand "crypto/rand"

	"github.com/keep-network/keep-core/config"
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
	The "genesis" subcommand triggers the first group selection. This action 
    can be done only once when there are no groups on the chain.`

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
				Usage:  "Performs genesis. Can be executed only one time.",
				Action: genesis,
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

	utility, err := ethereum.ConnectUtility(cfg.Ethereum)
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

	wait := make(chan struct{})

	utility.ThresholdRelay().OnSignatureSubmitted(func(entry *event.Entry) {
		requestMutex.Lock()
		defer requestMutex.Unlock()

		fmt.Fprintf(
			os.Stderr,
			"Relay entry received with value: [%v].\n",
			entry.Value,
		)

		wait <- struct{}{}
	})

	fmt.Printf("Requesting for a new relay entry at [%s]\n", time.Now())

	utility.RequestRelayEntry(seed).
		OnFailure(func(err error) {
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
	}
}

// genesis kicks off protocol to create the first group.
func genesis(c *cli.Context) error {
	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: [%v]", err)
	}

	utility, err := ethereum.ConnectUtility(cfg.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	err = utility.Genesis()
	if err != nil {
		return fmt.Errorf("error triggering genesis: [%v]", err)
	}
	return nil
}
