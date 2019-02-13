package local

import (
	"fmt"
	"math/big"
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

	monitor := NewStakeMonitor(big.NewInt(200))

	hasStake, err := monitor.HasMinimumStake("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitor.StakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitor.UnstakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)
}

func TestStakerFor(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	address := "0x65ea55c1f10491038425725dc00dffeab2a1e28a"
	staker, err := monitor.StakerFor(address)
	if err != nil {
		t.Fatal(err)
	}

	expectedStaker := &localStaker{
		address: address,
		stake:   big.NewInt(0),
	}
	if !reflect.DeepEqual(staker, expectedStaker) {
		t.Fatalf(
			"\nexpected: %+v\nactual:   %+v\n",
			expectedStaker,
			staker,
		)
	}

	stake, err := staker.Stake()
	if err != nil {
		t.Fatal(err)
	}
	if stake.Cmp(expectedStaker.stake) != 0 {
		t.Fatalf(
			"\nexpected: %v\nactual:   %v\n",
			big.NewInt(0),
			expectedStaker.stake,
		)
	}
}

func TestNoMinimumStakeByDefault(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	hasStake, err := monitor.HasMinimumStake("0x65ea55c1f10491038425725dc00dffeab2a1e28a")
	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("address should have no stake by default")
	}
}

func TestHasMinimumStakeIfStakedBefore(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	address := "0x524f2e0176350d950fa630d9a5a59a0a190daf48"

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
	monitor := NewStakeMonitor(big.NewInt(200))

	address := "0x524f2e0176350d950fa630d9a5a59a0a190daf48"

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

func TestStake(t *testing.T) {
	minimumStake := big.NewInt(200)
	expectedStake := new(big.Int).Mul(big.NewInt(5), minimumStake)

	monitor := NewStakeMonitor(minimumStake)
	address := "0x524f2e0176350d950fa630d9a5a59a0a190daf48"

	staker, err := monitor.StakerFor(address)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.StakeTokens(address)
	if err != nil {
		t.Fatal(err)
	}

	stake, err := staker.Stake()
	if err != nil {
		t.Fatal(err)
	}

	if stake.Cmp(expectedStake) != 0 {
		t.Fatalf(
			"\nexpected: %v\nactual:   %v\n",
			expectedStake,
			stake,
		)
	}
}
