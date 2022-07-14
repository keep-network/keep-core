package local

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
	localChain := Connect(testOperatorAddress)

	_, ok, err := localChain.OperatorToStakingProvider()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Expected operator not to be registered")
	}
}

func TestOperatorToStakingProvider(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)

	stakingProvider, ok, err := localChain.OperatorToStakingProvider()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Expected operator to be registered")
	}

	testutils.AssertStringsEqual(
		t,
		"staking provider",
		testStakingProviderAddress,
		stakingProvider.String(),
	)
}

func TestEligibleStake(t *testing.T) {
	localChain := Connect(testOperatorAddress)

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

	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(10))

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

func TestOperatorUpToDate_NotRegisteredOperator(t *testing.T) {
	localChain := Connect(testOperatorAddress)

	_, err := localChain.IsOperatorUpToDate()
	testutils.AssertErrorsEqual(t, errOperatorUnknown, err)
}

func TestOperatorUpToDate_NoStake(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if !isUpToDate {
		t.Fatal("expected the operator to be up to date")
	}
}

func TestOperatorUpToDate_ZeroStake(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(0))

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if !isUpToDate {
		t.Fatal("expected the operator to be up to date")
	}
}

func TestOperatorUpToDate_NonZeroStakeNotInPool(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(100))

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if isUpToDate {
		t.Fatal("expected the operator not to be up to date")
	}
}

func TestOperatorUpToDate_StakeInSyncWithWeight(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(100))
	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if !isUpToDate {
		t.Fatal("expected the operator to be up to date")
	}
}

func TestOperatorUpToDate_StakeNotInSyncWithWeight(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(100))
	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(101))

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if isUpToDate {
		t.Fatal("expected the operator not to be up to date")
	}
}

func TestJoinSortitionPool_NotRegisteredOperator(t *testing.T) {
	localChain := Connect(testOperatorAddress)

	err := localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, errOperatorUnknown, err)
}

func TestJoinSortitionPool_AuthorizationBelowMinimum(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)

	err := localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, errAuthorizationBelowMinimum, err)
}

func TestJoinSortitionPool(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(1))

	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}
}

func TestJoinSortitionPool_OperatorAlreadyInPool(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(1))

	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}

	err = localChain.JoinSortitionPool()
	testutils.AssertErrorsEqual(t, errOperatorAlreadyRegisteredInPool, err)
}

func TestUpdateOperatorStatus_NotRegisteredOperator(t *testing.T) {
	localChain := Connect(testOperatorAddress)

	err := localChain.UpdateOperatorStatus()
	testutils.AssertErrorsEqual(t, errOperatorUnknown, err)
}

func TestUpdateOperatorStatus(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.RegisterOperator(testStakingProviderAddress, testOperatorAddress)
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(100))
	err := localChain.JoinSortitionPool()
	if err != nil {
		t.Fatal(err)
	}
	localChain.SetEligibleStake(testStakingProviderAddress, big.NewInt(101))

	isUpToDate, err := localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if isUpToDate {
		t.Fatal("expected the operator not to be up to date")
	}

	localChain.UpdateOperatorStatus()

	isUpToDate, err = localChain.IsOperatorUpToDate()
	if err != nil {
		t.Fatal(err)
	}

	if !isUpToDate {
		t.Fatal("expected the operator to be up to date")
	}
}

func TestIsEligibleForRewards_EligibleOperator(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(0))

	isEligibileForRewards, err := localChain.IsEligibleForRewards()
	if err != nil {
		t.Fatal(err)
	}

	if !isEligibileForRewards {
		t.Fatal("expected the operator to be eligible for rewards")
	}
}

func TestIsEligibleForRewards_NotEligibleOperator(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(1))

	isEligibileForRewards, err := localChain.IsEligibleForRewards()
	if err != nil {
		t.Fatal(err)
	}

	if isEligibileForRewards {
		t.Fatal("expected the operator not to be eligible for rewards")
	}
}

func TestCanRestoreRewardEligibility_Eligible(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(0))
	localChain.SetCurrentTimestamp(big.NewInt(1))

	canRestoreRewardEligibility, err := localChain.CanRestoreRewardEligibility()
	if err != nil {
		t.Fatal(err)
	}

	if !canRestoreRewardEligibility {
		t.Fatal("expected the operator can restore reward eligibility")
	}
}

func TestCanRestoreRewardEligibility_NotEligible(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(1))
	localChain.SetCurrentTimestamp(big.NewInt(1))

	canRestoreRewardEligibility, err := localChain.CanRestoreRewardEligibility()
	if err != nil {
		t.Fatal(err)
	}

	if canRestoreRewardEligibility {
		t.Fatal("expected the operator cannot restore reward eligibility")
	}
}

func TestRestoreRewardEligibility_Restore(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(1))
	localChain.SetCurrentTimestamp(big.NewInt(2))

	isEligibileForRewards, err := localChain.IsEligibleForRewards()
	if err != nil {
		t.Fatal(err)
	}

	if isEligibileForRewards {
		t.Fatal("expected the operator not to be eligible for rewards")
	}

	err = localChain.RestoreRewardEligibility()
	if err != nil {
		t.Fatal(err)
	}

	isEligibileForRewards, err = localChain.IsEligibleForRewards()
	if err != nil {
		t.Fatal(err)
	}

	if !isEligibileForRewards {
		t.Fatal("expected the operator to be eligible for rewards")
	}
}

func TestRestoreRewardEligibility_CannotRestore(t *testing.T) {
	localChain := Connect(testOperatorAddress)
	localChain.SetRewardIneligibility(big.NewInt(2))
	localChain.SetCurrentTimestamp(big.NewInt(1))

	isEligibileForRewards, err := localChain.IsEligibleForRewards()
	if err != nil {
		t.Fatal(err)
	}

	if isEligibileForRewards {
		t.Fatal("expected the operator to not be eligible for rewards")
	}

	err = localChain.RestoreRewardEligibility()
	if err == nil {
		t.Fatal("expected the operator not to be eligible for rewards")
	}
}
