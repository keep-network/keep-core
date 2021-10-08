import { ContractsLoaded, getContractDeploymentBlockNumber } from "../contracts"
import { TOKEN_STAKING_ESCROW_CONTRACT_NAME } from "../constants/constants"

export const fetchEscrowDepositsByGrantId = async (grantId) => {
  const { tokenStakingEscrow } = await ContractsLoaded

  return await tokenStakingEscrow.getPastEvents("Deposited", {
    fromBlock: await getContractDeploymentBlockNumber(
      TOKEN_STAKING_ESCROW_CONTRACT_NAME
    ),
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
