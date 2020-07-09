import { contractService } from "./contracts.service"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import { getKeepRandomBeaconOperatorAddress } from "../contracts"
import { getOperatorsOfAuthorizer } from "./token-staking.service"

const keepRandomBeaconOperatorAddress = getKeepRandomBeaconOperatorAddress()

const fetchRandomBeaconAuthorizationData = async (web3Context) => {
  const { yourAddress } = web3Context

  const authorizerOperators = await getOperatorsOfAuthorizer(
    web3Context,
    yourAddress
  )
  const authorizationData = []
  // Fetch all authorizer operators
  for (let i = 0; i < authorizerOperators.length; i++) {
    const operatorAddress = authorizerOperators[i]

    const delegatedTokens = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "getDelegationInfo",
      operatorAddress
    )

    const isKeepRandomBeaconOperatorAuthorized = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "isAuthorizedForOperator",
      operatorAddress,
      keepRandomBeaconOperatorAddress
    )

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
