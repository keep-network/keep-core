package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain"
)

// RolesOf returns the stake owner, the beneficiary and the authorizer for the
// specified staking provider address. If the owner is set, the function considers
// the staking provider to have a stake delegation and returns the boolean flag
// set to true.
func (ec *Chain) RolesOf(stakingProvider chain.Address) (
	owner, beneficiary, authorizer chain.Address, hasStake bool, err error,
) {
	rolesOf, err := ec.tokenStaking.RolesOf(
		common.HexToAddress(stakingProvider.String()),
	)

	if err != nil {
		return "", "", "", false, err
	}

	if (rolesOf.Owner == common.Address{}) {
		return "", "", "", false, nil
	}

	return chain.Address(rolesOf.Owner.Hex()),
		chain.Address(rolesOf.Beneficiary.Hex()),
		chain.Address(rolesOf.Authorizer.Hex()),
		true,
		nil
}
