import {
  getKeepRandomBeaconOperatorAddress,
  ContractsLoaded,
} from "../contracts"
import { getOperatorsOfAuthorizer } from "./token-staking.service"

const keepRandomBeaconOperatorAddress = getKeepRandomBeaconOperatorAddress()

const fetchRandomBeaconAuthorizationData = async (address) => {
  const { stakingContract } = await ContractsLoaded
  const authorizerOperators = await getOperatorsOfAuthorizer(address)
  const authorizationData = []
  // Fetch all authorizer operators
  for (let i = 0; i < authorizerOperators.length; i++) {
    const operatorAddress = authorizerOperators[i]

    const delegatedTokens = await stakingContract.methods
      .getDelegationInfo(operatorAddress)
      .call()

    const isKeepRandomBeaconOperatorAuthorized = await stakingContract.methods
      .isAuthorizedForOperator(operatorAddress, keepRandomBeaconOperatorAddress)
      .call()

    const authorizerOperator = {
      operatorAddress: operatorAddress,
      stakeAmount: delegatedTokens.amount,
      contracts: [
        {
          contractName: "Keep Random Beacon Operator Contract",
          operatorContractAddress: keepRandomBeaconOperatorAddress,
          isAuthorized: isKeepRandomBeaconOperatorAuthorized,
        },
      ],
    }

    authorizationData.push(authorizerOperator)
  }

  return authorizationData
}

const authorizeKeepRandomBeaconOperatorContract = async (
  web3Context,
  operatorAddress,
  onTransactionHashCallback
) => {
  const { stakingContract, yourAddress } = web3Context
  try {
    await stakingContract.methods
      .authorizeOperatorContract(
        operatorAddress,
        keepRandomBeaconOperatorAddress
      )
      .send({ from: yourAddress })
      .on("transactionHash", onTransactionHashCallback)
  } catch (error) {
    throw error
  }
}

export const beaconAuthorizationService = {
  fetchRandomBeaconAuthorizationData,
  authorizeKeepRandomBeaconOperatorContract,
}
