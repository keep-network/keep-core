import {
  ContractsLoaded,
  getContractDeploymentBlockNumber,
  getThresholdTokenStakingAddress,
} from "../contracts"
import { getOperatorsOfAuthorizer } from "./token-staking.service"
import {
  AUTH_CONTRACTS_LABEL,
  TOKEN_STAKING_CONTRACT_NAME,
} from "../constants/constants"
import { Keep } from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"

const fetchThresholdAuthorizationData = async (address) => {
  if (!address) {
    return []
  }
  const thresholdTokenStakingContractAddress = getThresholdTokenStakingAddress()
  const { stakingContract } = await ContractsLoaded
  const authorizerOperators = await getOperatorsOfAuthorizer(address)
  const authorizationData = []

  const keepToTStakedEvents =
    await Keep.keepToTStaking.getStakedEventsByOperator(address)

  const operatorsStaked = keepToTStakedEvents.map(
    (event) => event.returnValues.stakingProvider
  )

  const stakesParticipants = (
    await stakingContract.getPastEvents("OperatorStaked", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_CONTRACT_NAME
      ),
      filter: { operator: authorizerOperators },
    })
  ).map((event) => {
    return {
      authorizer: event.returnValues.authorizer,
      operator: event.returnValues.operator,
      beneficiary: event.returnValues.beneficiary,
    }
  })

  // Fetch all authorizer operators
  for (let i = 0; i < authorizerOperators.length; i++) {
    const operatorAddress = authorizerOperators[i]

    const stakeParticipant = stakesParticipants.find((participants) => {
      return isSameEthAddress(participants.authorizer, address)
    })

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
      authorizerAddress: stakeParticipant.authorizer,
      operatorAddress: operatorAddress,
      beneficiaryAddress: stakeParticipant.beneficiary,
      stakeAmount: delegatedTokens.amount,
      contracts: [
        {
          contractName: AUTH_CONTRACTS_LABEL.THRESHOLD_TOKEN_STAKING,
          operatorContractAddress: thresholdTokenStakingContractAddress,
          isAuthorized: isThresholdTokenStakingContractAuthorized,
        },
      ],
      isStakedToT: operatorsStaked.some((operatorStaked) =>
        isSameEthAddress(operatorStaked, operatorAddress)
      ),
    }

    authorizationData.push(authorizerOperator)
  }

  return authorizationData
}

export const thresholdAuthorizationService = {
  fetchThresholdAuthorizationData,
}
