package local

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDetectInvalidAddress(t *testing.T) {
	assertInvalidEthereumAddress := func(hasStake bool, err error, t *testing.T) {
		expectedError := fmt.Errorf("not a valid ethereum address: 0x010102003")

		if !reflect.DeepEqual(expectedError, err) {
			t.Fatalf(
				"unexpected error\nexpected: %v\nactual: %v",
				expectedError,
				err,
			)
		}
		if hasStake {
			t.Fatalf("expected 'false' result")
		}
	}

	monitor := NewStakeMonitor()

	hasStake, err := monitor.HasMinimumStake("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitor.StakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitor.UnstakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)
}

func TestNoMinimumStakeByDefault(t *testing.T) {
	monitor := NewStakeMonitor()

	hasStake, err := monitor.HasMinimumStake(
		"0x65ea55c1f10491038425725dc00dffeab2a1e28a",
	)

	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("address should have no stake by default")
	}
}

func TestHasMinimumStakeIfStakedBefore(t *testing.T) {
	monitor := NewStakeMonitor()

	address := "0x524f2e0176350d950fa630d9a5a59a0a190daf48"
	monitor.StakerFor(address)

	err := monitor.StakeTokens(address)
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitor.HasMinimumStake(address)

	if err != nil {
		t.Fatal(err)
	}

	if !hasStake {
		t.Fatal("address should have tokens staked")
	}
}

func TestNoMinimumStakeIfUnstaked(t *testing.T) {
	monitor := NewStakeMonitor()

	address := "0x524f2e0176350d950fa630d9a5a59a0a190daf48"
	monitor.StakerFor(address)

	err := monitor.StakeTokens(address)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.UnstakeTokens(address)
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitor.HasMinimumStake(address)
	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("address should have no stake if unstaked earlier")
	}
}
