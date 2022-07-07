package local

import (
	"github.com/keep-network/keep-core/pkg/operator"
	"math/big"
	"reflect"
	"testing"
)

func TestStakerFor(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	staker, err := monitor.stakerFor(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	expectedStaker := &localStaker{
		publicKey: operatorPublicKey,
		stake:     big.NewInt(0),
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

	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitor.HasMinimumStake(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("operator should have no stake by default")
	}
}

func TestHasMinimumStakeIfStakedBefore(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.StakeTokens(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitor.HasMinimumStake(operatorPublicKey)

	if err != nil {
		t.Fatal(err)
	}

	if !hasStake {
		t.Fatal("operator should have tokens staked")
	}
}

func TestNoMinimumStakeIfUnstaked(t *testing.T) {
	monitor := NewStakeMonitor(big.NewInt(200))

	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.StakeTokens(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.UnstakeTokens(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitor.HasMinimumStake(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("operator should have no stake if unstaked earlier")
	}
}

func TestStake(t *testing.T) {
	minimumStake := big.NewInt(200)
	expectedStake := new(big.Int).Mul(big.NewInt(5), minimumStake)

	monitor := NewStakeMonitor(minimumStake)

	_, operatorPublicKey, err := operator.GenerateKeyPair(DefaultCurve)
	if err != nil {
		t.Fatal(err)
	}

	staker, err := monitor.stakerFor(operatorPublicKey)
	if err != nil {
		t.Fatal(err)
	}

	err = monitor.StakeTokens(operatorPublicKey)
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
