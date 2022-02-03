import {
  ContractsLoaded,
  getContractDeploymentBlockNumber,
  getThresholdTokenStakingAddress,
} from "../contracts"
import { getAllOperatorStakedEventsByAuthorizer } from "./token-staking.service"
import {
  AUTH_CONTRACTS_LABEL,
  TOKEN_GRANT_CONTRACT_NAME,
} from "../constants/constants"
import { Keep } from "../contracts"
import { isSameEthAddress } from "../utils/general.utils"
import { gt } from "../utils/arithmetics.utils"

const fetchThresholdAuthorizationData = async (address) => {
  if (!address) {
    return []
  }
  const thresholdTokenStakingContractAddress = getThresholdTokenStakingAddress()
  const { stakingContract, grantContract } = await ContractsLoaded
  const keepOperatorStakedEvents = await getAllOperatorStakedEventsByAuthorizer(
    address
  )
  const authorizerOperators = keepOperatorStakedEvents.map(
    (_) => _.returnValues.operator
  )
  const authorizationData = []

  const keepToTStakedEvents =
    await Keep.keepToTStaking.getStakedEventsByOperator(authorizerOperators)

  const operatorsStakedToT = keepToTStakedEvents.map(
    (event) => event.returnValues.stakingProvider
  )

  const tokenGrantStakingEvents = (
    await grantContract.getPastEvents("TokenGrantStaked", {
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_GRANT_CONTRACT_NAME
      ),
      filter: { operator: authorizerOperators },
    })
  ).map((event) => {
    return {
      operator: event.returnValues.operator,
    }
  })

  // Fetch all authorizer operators
  for (let i = 0; i < authorizerOperators.length; i++) {
    const operatorAddress = authorizerOperators[i]

    const stakeParticipant = keepOperatorStakedEvents.find((event) => {
      return isSameEthAddress(operatorAddress, event.returnValues.operator)
    })

    const { amount: stakeAmount, undelegatedAt } = await stakingContract.methods
      .getDelegationInfo(operatorAddress)
      .call()

    // If stake is undelegated we won't display it, because undelegated stakes
    // can't be staked to Threshold
    if (undelegatedAt !== "0" && gt(stakeAmount, 0)) continue

    const isThresholdTokenStakingContractAuthorized =
      await stakingContract.methods
        .isAuthorizedForOperator(
          operatorAddress,
          thresholdTokenStakingContractAddress
        )
        .call()

    const authorizerOperator = {
      authorizerAddress: stakeParticipant.returnValues.authorizer,
      operatorAddress: operatorAddress,
      beneficiaryAddress: stakeParticipant.returnValues.beneficiary,
      stakeAmount: stakeAmount,
      contracts: [
        {
          contractName: AUTH_CONTRACTS_LABEL.THRESHOLD_TOKEN_STAKING,
          operatorContractAddress: thresholdTokenStakingContractAddress,
          isAuthorized: isThresholdTokenStakingContractAuthorized,
        },
      ],
      isStakedToT: operatorsStakedToT.some((operatorStaked) =>
        isSameEthAddress(operatorStaked, operatorAddress)
      ),
      isFromGrant: tokenGrantStakingEvents.some((tokenGrantStakingEvent) =>
        isSameEthAddress(tokenGrantStakingEvent.operator, operatorAddress)
      ),
    }

    authorizationData.push(authorizerOperator)
  }

  return authorizationData
}

export const thresholdAuthorizationService = {
  fetchThresholdAuthorizationData,
}
