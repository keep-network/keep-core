package local

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDetectInvalidAddress(t *testing.T) {
	monitoring := newLocalStakeMonitoring()

	hasStake, err := monitoring.HasMinimumStake("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitoring.stakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)

	err = monitoring.unstakeTokens("0x010102003")
	assertInvalidEthereumAddress(hasStake, err, t)
}

func assertInvalidEthereumAddress(hasStake bool, err error, t *testing.T) {
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

func TestNoMinimumStakeByDefault(t *testing.T) {
	monitoring := newLocalStakeMonitoring()

	hasStake, err := monitoring.HasMinimumStake(
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
	monitoring := newLocalStakeMonitoring()

	err := monitoring.stakeTokens("0x524f2e0176350d950fa630d9a5a59a0a190daf48")
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitoring.HasMinimumStake(
		"0x524f2e0176350d950fa630d9a5a59a0a190daf48",
	)

	if err != nil {
		t.Fatal(err)
	}

	if !hasStake {
		t.Fatal("address should have tokens staked")
	}
}

func TestNoMinimumStakeIfUnstaked(t *testing.T) {
	monitoring := newLocalStakeMonitoring()

	err := monitoring.stakeTokens("0x524f2e0176350d950fa630d9a5a59a0a190daf48")
	if err != nil {
		t.Fatal(err)
	}

	err = monitoring.unstakeTokens("0x524f2e0176350d950fa630d9a5a59a0a190daf48")
	if err != nil {
		t.Fatal(err)
	}

	hasStake, err := monitoring.HasMinimumStake(
		"0x524f2e0176350d950fa630d9a5a59a0a190daf48",
	)

	if err != nil {
		t.Fatal(err)
	}

	if hasStake {
		t.Fatal("address should have no stake if unstaked earlier")
	}
}
