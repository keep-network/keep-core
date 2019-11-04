package groupselection

import (
	"math/big"
	"testing"
)

var stakingAddress = []byte("staking address")
var previousBeaconOutput = []byte("test beacon output")

func testNaturalThreshold() *big.Int { // 2^256 / 2
	return new(big.Int).Div(
		new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(2),
	)
}

func TestAllTicketsGenerated(t *testing.T) {
	minimumStake := big.NewInt(20)
	availableStake := big.NewInt(1000)
	virtualStakers := availableStake.Int64() / minimumStake.Int64()

	initialTickets, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		testNaturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// We should have 1000/20 = 50 tickets
	allTicketsCount := len(initialTickets) + len(reactiveTickets)
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

	initialTickets, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		testNaturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	allTickets := append(initialTickets, reactiveTickets...)

	// Tickets should be sorted in ascending order
	for i := 0; i < len(allTickets)-1; i++ {
		value := allTickets[i].intValue()
		nextValue := allTickets[i+1].intValue()

		if value.Cmp(nextValue) > 0 {
			t.Errorf("tickets not sorted in ascending order")
		}
	}
}

func TestInitialTicketsGeneatedBelowNaturalThreshold(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10000)

	initialTickets, _, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		testNaturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// All initial submission tickets should have value below the natural
	// threshold
	for _, ticket := range initialTickets {
		if ticket.intValue().Cmp(testNaturalThreshold()) >= 0 {
			t.Errorf(
				"initial submission ticket value should be below natural "+
					"threshold\nvalue:     [%v]\nthreshold: [%v]",
				ticket.intValue(),
				testNaturalThreshold(),
			)
		}
	}
}

func TestReactiveTicketsGeneatedAboveNaturalThreshold(t *testing.T) {
	minimumStake := big.NewInt(1)
	availableStake := big.NewInt(10000)

	_, reactiveTickets, err := generateTickets(
		previousBeaconOutput,
		stakingAddress,
		availableStake,
		minimumStake,
		testNaturalThreshold(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// All reactive submission tickets should have value equal or above
	// the natural threshold
	for _, ticket := range reactiveTickets {
		if ticket.intValue().Cmp(testNaturalThreshold()) <= 0 {
			t.Errorf(
				"reactive submission ticket value should not be below natural "+
					"threshold\nvalue:     [%v]\nthreshold: [%v]",
				ticket.intValue(),
				testNaturalThreshold(),
			)
		}
	}
}
