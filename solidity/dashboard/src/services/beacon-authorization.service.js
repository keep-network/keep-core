import { contractService } from "./contracts.service"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { 
  CONTRACT_DEPLOY_BLOCK_NUMBER, 
  getKeepRandomBeaconOperatorAddress,
} from "../contracts"

const keepRandomBeaconOperatorAddress = getKeepRandomBeaconOperatorAddress()

const fetchStakedEvents = async (web3Context) => {
  return contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "Staked",
    { fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME] }
  )
}

const fetchRandomBeaconAuthorizationData = async (web3Context) => {
  const { yourAddress } = web3Context
  const stakedEvents = await fetchStakedEvents(web3Context)

  const visitedOperators = {}
  const authorizerOperators = []

  // Fetch all authorizer operators
  for (let i = 0; i < stakedEvents.length; i++) {
    const {
      returnValues: { from: operatorAddress },
    } = stakedEvents[i]

    if (visitedOperators.hasOwnProperty(operatorAddress)) {
      continue
    }
    visitedOperators[operatorAddress] = operatorAddress
    const authorizerOfOperator = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "authorizerOf",
      operatorAddress
    )
    if (isSameEthAddress(authorizerOfOperator, yourAddress)) {
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

      authorizerOperators.push(authorizerOperator)
    }
  }

  return authorizerOperators
}

export const beaconAuthorizationService = {
  fetchRandomBeaconAuthorizationData,
}
