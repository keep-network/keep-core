package groupselection

import (
	"math/big"
	"testing"
)

var stakingAddress = []byte("staking address")
var previousBeaconOutput = []byte("test beacon output")

func TestAllTicketsGenerated(t *testing.T) {
	minimumStake := big.NewInt(20)
	availableStake := big.NewInt(1000)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	tickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
	)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 1000/20 = 50 tickets
	allTicketsCount := len(tickets)
	if allTicketsCount != int(virtualStakers) {
		t.Fatalf(
			"expected [%d] tickets, has [%d] tickets",
			virtualStakers,
			allTicketsCount,
		)
	}
}

func TestTicketsGeneratedInOrder(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(100)

	tickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Tickets should be sorted in ascending order
	for i := 0; i < len(tickets)-1; i++ {
		value := tickets[i].intValue()
		nextValue := tickets[i+1].intValue()

		if value.Cmp(nextValue) > 0 {
			t.Errorf("tickets not sorted in ascending order")
		}
	}
}
