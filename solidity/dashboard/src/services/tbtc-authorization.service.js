import { contractService } from "./contracts.service"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import { BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME } from "../constants/constants"
import { KEEP_BONDING_CONTRACT_NAME } from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import { 
    CONTRACT_DEPLOY_BLOCK_NUMBER,
    getBondedECDSAKeepFactoryAddress,
    getTBTCSystemAddress,
    getKeepRandomBeaconOperatorAddress,
} from "../contracts"

const tBTCSystemAddress = getTBTCSystemAddress()
const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()

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

  console.log("getBondedECDSAKeepFactoryAddress: ", getBondedECDSAKeepFactoryAddress())
  console.log("getTBTCSystemAddress: ", getTBTCSystemAddress())
  console.log("getKeepRandomBeaconOperatorAddress: ", getKeepRandomBeaconOperatorAddress())
  
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

const authorizeBondedECDSAKeepFactory = async (web3Context, operatorAddress) => {
    await contractService.makeCall(
        web3Context,
        TOKEN_STAKING_CONTRACT_NAME,
        "authorizeOperatorContract",
        operatorAddress,
        bondedECDSAKeepFactoryAddress
    )
}

const authorizeTBTCSystem = async (web3Context, operatorAddress) => {
  try {
    const sortitionPoolAddress = await contractService.makeCall(
      web3Context,
      BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
      "getSortitionPool",
      tBTCSystemAddress,
    )
    
    await contractService.makeCall(
      web3Context,
      KEEP_BONDING_CONTRACT_NAME,
      "authorizeSortitionPoolContract",
      operatorAddress,
      sortitionPoolAddress,
    )
  } catch(error) {
    // TODO: handle the error properly
    console.error("failed to authorize tBTC application", error)
  }
}

export const tbtcAuthorizationService = {
    fetchTBTCAuthorizationData,
    authorizeBondedECDSAKeepFactory,
    authorizeTBTCSystem,
}