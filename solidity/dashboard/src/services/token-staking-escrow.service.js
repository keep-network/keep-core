import { ContractsLoaded, CONTRACT_DEPLOY_BLOCK_NUMBER } from "../contracts"

export const fetchEscrowDepositsByGrantId = async (grantId) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  return await tokenStakingEscrow.getPastEvents("Deposited", {
    fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
    filter: { grantId },
  })
}

export const fetchWithdrawableAmountForDeposit = async (operatorAddress) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  return await tokenStakingEscrow.methods.withdrawable(operatorAddress).call()
}

export const fetchDepositWithdrawnAmount = async (operatorAddress) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  return await tokenStakingEscrow.methods
    .depositWithdrawnAmount(operatorAddress)
    .call()
}

export const fetchDepositAvailableAmount = async (operatorAddress) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  return await tokenStakingEscrow.methods
    .availableAmount(operatorAddress)
    .call()
}
