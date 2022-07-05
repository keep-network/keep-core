package sortition

import (
	"math/big"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

const (
	testOperatorAddress        = "0x3c5eBAcFe5aE12D82d43602a12b8bBb76b893CfA"
	testStakingProviderAddress = "0x80C63B577DC79B2432357BECC5b431dfb8E181DD"
	testThirdPartyAddress      = "0x91605Ef3251fb8bd5e12Cad7F897f1e0c2183Ceb"
)

func TestOperatorToStakingProvider_NotRegisteredOperator(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)

	_, err := localChain.OperatorToStakingProvider()

	testutils.AssertErrorsEqual(t, ErrOperatorNotRegistered, err)
}

func TestOperatorToStakingProvider(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)
	localChain.registerOperator(testStakingProviderAddress, testOperatorAddress)

	stakingProvider, err := localChain.OperatorToStakingProvider()
	if err != nil {
		t.Fatal(err)
	}

	testutils.AssertStringsEqual(
		t,
		"staking provider",
		testStakingProviderAddress,
		stakingProvider.String(),
	)
}

func TestEligibleStake(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)

	eligibleStake, err := localChain.EligibleStake(testStakingProviderAddress)
	if err != nil {
		t.Fatal(err)
	}
	testutils.AssertBigIntsEqual(
		t,
		"eligible stake",
		eligibleStake,
		big.NewInt(0),
	)

	localChain.setEligibleStake(testStakingProviderAddress, big.NewInt(10))

	eligibleStake, err = localChain.EligibleStake(testStakingProviderAddress)
	if err != nil {
		t.Fatal(err)
	}
	testutils.AssertBigIntsEqual(
		t,
		"eligible stake",
		eligibleStake,
		big.NewInt(10),
	)
}

func TestJoinSortitionPool_NotRegisteredOperator(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)

	err := localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, ErrOperatorNotRegistered, err)
}

func TestJoinSortitionPool_AuthorizationBelowMinimum(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)
	localChain.registerOperator(testStakingProviderAddress, testOperatorAddress)

	err := localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, errAuthorizationBelowMinimum, err)
}

func TestJoinSortitionPool(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)
	localChain.registerOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.setEligibleStake(testStakingProviderAddress, big.NewInt(1))

	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}
}

func TestJoinSortitionPool_OperatorAlreadyInPool(t *testing.T) {
	localChain := connectLocal(testOperatorAddress)
	localChain.registerOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.setEligibleStake(testStakingProviderAddress, big.NewInt(1))

	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}

	err = localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, errOperatorAlreadyRegisteredInPool, err)
}
