package cmd

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/urfave/cli"
)

// RelayRequest requests a new entry from the threshold relay and prints the
// request id. By default, it also waits until the associated relay entry is
// generated and prints out the entry.
func RelayRequest(c *cli.Context) error {
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
		fmt.Println(
			"Relay entry request submitted with id ",
			request.RequestID.String(),
			".",
		)
		requestMutex.Lock()
		requestID = request.RequestID
		requestMutex.Unlock()
	})

	provider.ThresholdRelay().OnRelayEntryGenerated(func(entry *event.Entry) {
		requestMutex.Lock()
		defer requestMutex.Unlock()

		if requestID != nil && requestID.Cmp(entry.RequestID) == 0 {
			valueBigInt := &big.Int{}
			valueBigInt.SetBytes(entry.Value[:])
			fmt.Printf(
				"Relay entry received with value: [%s].\n",
				valueBigInt.String(),
			)

			wait <- struct{}{}
		}
	})

	// provider.ThresholdRelay().RequestRelayEntry(&big.Int{}, &big.Int{})

	select {
	case <-wait:
		cancel()
		os.Exit(0)
	case <-ctx.Done():
		err := ctx.Err()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Request errored out: [%v].\n",
				err,
			)
		} else {
			fmt.Fprintf(os.Stderr, "Request errored for unknown reason.\n")
		}

		os.Exit(1)
	}

	return nil
}

// RelayEntry requests an entry with a particular id from the threshold relay
// and prints that entry.
func RelayEntry(c *cli.Context) error {
	return fmt.Errorf("relay entry lookups are currently unimplemented")
}
