import { ContractsLoaded, getThresholdTokenStakingAddress } from "../contracts"
import { getOperatorsOfAuthorizer } from "./token-staking.service"
import { AUTH_CONTRACTS_LABEL } from "../constants/constants"

const fetchThresholdAuthorizationData = async (address) => {
  if (!address) {
    return []
  }
  const thresholdTokenStakingContractAddress = getThresholdTokenStakingAddress()
  const { stakingContract } = await ContractsLoaded
  const authorizerOperators = await getOperatorsOfAuthorizer(address)
  const authorizationData = []
  // Fetch all authorizer operators
  for (let i = 0; i < authorizerOperators.length; i++) {
    const operatorAddress = authorizerOperators[i]

    const delegatedTokens = await stakingContract.methods
      .getDelegationInfo(operatorAddress)
      .call()

    const isThresholdTokenStakingContractAuthorized =
      await stakingContract.methods
        .isAuthorizedForOperator(
          operatorAddress,
          thresholdTokenStakingContractAddress
        )
        .call()

    const authorizerOperator = {
      operatorAddress: operatorAddress,
      stakeAmount: delegatedTokens.amount,
      contracts: [
        {
          contractName: AUTH_CONTRACTS_LABEL.THRESHOLD_TOKEN_STAKING,
          operatorContractAddress: thresholdTokenStakingContractAddress,
          isAuthorized: isThresholdTokenStakingContractAuthorized,
        },
      ],
    }

    authorizationData.push(authorizerOperator)
  }

  return authorizationData
}

export const thresholdAuthorizationService = {
  fetchThresholdAuthorizationData,
}
