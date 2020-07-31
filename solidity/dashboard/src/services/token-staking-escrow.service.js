import { contracts, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"

export const fetchEscrowDepositsByGrantId = async (grantId) => {
  return await contracts.tokenStakingEscrow.getPastEvents("Deposited", {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
    filter: { grantId },
  })
}

export const fetchWithdrawableAmountForDeposit = async (operatorAddress) => {
  return await contracts.tokenStakingEscrow.methods
    .withdrawable(operatorAddress)
    .call()
}

export const fetchDepositWithdrawnAmount = async (operatorAddress) => {
  return await contracts.tokenStakingEscrow.methods
    .depositWithdrawnAmount(operatorAddress)
    .call()
}
