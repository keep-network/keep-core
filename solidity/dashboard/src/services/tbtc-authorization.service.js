import { contractService } from "./contracts.service"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { 
    CONTRACT_DEPLOY_BLOCK_NUMBER,
    getBondedECDSAKeepFactoryAddress,
    getTBTCSystemAddress
} from "../contracts"

const fetchTBTCAuthorizationData = async (web3Context) => {
  const { yourAddress } = web3Context

  const stakedEvents = await contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "Staked",
    { fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME] }
  )
  const visitedOperators = {}
  const authorizerOperators = []
  const data = {}

  console.log("stakedEvents.length: ", stakedEvents.length)
  console.log("getBondedECDSAKeepFactoryAddress: ", getBondedECDSAKeepFactoryAddress())
  console.log("getTBTCSystemAddress: ", getTBTCSystemAddress())

  // TODO: need to double check this
    const tBTCSystemAddress = getTBTCSystemAddress()
    const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()
  
    // Fetch all authorizer operators
  for (let i = 0; i < stakedEvents.length; i++) {
    const {
      returnValues: { from: operatorAddress },
    } = stakedEvents[i]

    if (visitedOperators.hasOwnProperty(operatorAddress)) {
      continue
    }
    visitedOperators[operatorAddress] = operatorAddress
    console.log("operatorAddress: ", operatorAddress)
    const authorizerOfOperator = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "authorizerOf",
      operatorAddress
    )
    console.log("authorizerOfOperator: ", authorizerOfOperator)
    if (isSameEthAddress(authorizerOfOperator, yourAddress)) {
        
        const delegatedTokens = await contractService.makeCall(
            web3Context,
            TOKEN_STAKING_CONTRACT_NAME,
            "getDelegationInfo",
            operatorAddress
        )

        const isBondedECDSAKeepFactoryAuthorized = await contractService.makeCall(
            web3Context,
            TOKEN_STAKING_CONTRACT_NAME,
            "isAuthorizedForOperator",
            operatorAddress,
            bondedECDSAKeepFactoryAddress
        )

        const isTBTCSystemAuthorized = await contractService.makeCall(
            web3Context,
            TOKEN_STAKING_CONTRACT_NAME,
            "isAuthorizedForOperator",
            operatorAddress,
            tBTCSystemAddress
         )

        const authorizerOperator = {
            operatorAddress: operatorAddress,
            stakeAmount: delegatedTokens.amount,
            contracts: [
                {
                    contractName: "BondedECDSAKeepFactory",
                    operatorContractAddress: bondedECDSAKeepFactoryAddress,
                    isAuthorized: isBondedECDSAKeepFactoryAuthorized,
                },
                {
                    contractName: "TBTCSystem",
                    operatorContractAddress: tBTCSystemAddress,
                    isAuthorized: isTBTCSystemAuthorized,
                },
            ]
        }

        authorizerOperators.push(authorizerOperator)

    }
  }

  return authorizerOperators

}

export const tbtcAuthorizationService = {
    fetchTBTCAuthorizationData,
}