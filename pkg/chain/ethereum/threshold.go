package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/pkg/chain"
)

func (ec *Chain) RolesOf(stakingProvider chain.Address) (
	owner, beneficiary, authorizer chain.Address, ok bool, err error,
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
