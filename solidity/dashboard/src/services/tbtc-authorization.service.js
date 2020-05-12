import { contractService } from "./contracts.service"
import { TOKEN_STAKING_CONTRACT_NAME } from "../constants/constants"
import {
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  getBondedECDSAKeepFactoryAddress,
  getTBTCSystemAddress,
  getKeepRandomBeaconOperatorAddress,
} from "../contracts"
import web3Utils from "web3-utils"

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

  console.log(
    "getBondedECDSAKeepFactoryAddress: ",
    getBondedECDSAKeepFactoryAddress()
  )
  console.log("getTBTCSystemAddress: ", getTBTCSystemAddress())
  console.log(
    "getKeepRandomBeaconOperatorAddress: ",
    getKeepRandomBeaconOperatorAddress()
  )

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

      const isTBTCSystemAuthorized = await isTbtcSystemAuthorized(
        web3Context,
        operatorAddress
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
        ],
      }

      authorizerOperators.push(authorizerOperator)
    }
  }

  return authorizerOperators
}

const isTbtcSystemAuthorized = async (web3Context, operatorAddress) => {
  try {
    const sortitionPoolAddress = await contractService.makeCall(
      web3Context,
      BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
      "getSortitionPool",
      tBTCSystemAddress
    )

    return await contractService.makeCall(
      web3Context,
      KEEP_BONDING_CONTRACT_NAME,
      "hasSecondaryAuthorization",
      operatorAddress,
      sortitionPoolAddress
    )
  } catch {
    return false
  }
}

const authorizeBondedECDSAKeepFactory = async (
  web3Context,
  operatorAddress,
  onTransactionHashCallback
) => {
  const { stakingContract, yourAddress } = web3Context
  try {
    await stakingContract.methods
      .authorizeOperatorContract(operatorAddress, bondedECDSAKeepFactoryAddress)
      .send({ from: yourAddress })
      .on("transactionHash", onTransactionHashCallback)
  } catch (error) {
    throw error
  }
}

const authorizeTBTCSystem = async (
  web3Context,
  operatorAddress,
  onTransactionHashCallback
) => {
  const { keepBondingContract, yourAddress } = web3Context
  try {
    const sortitionPoolAddress = await contractService.makeCall(
      web3Context,
      BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
      "getSortitionPool",
      tBTCSystemAddress
    )

    await keepBondingContract.methods
      .authorizeSortitionPoolContract(operatorAddress, sortitionPoolAddress)
      .send({ from: yourAddress })
      .on("transactionHash", onTransactionHashCallback)
  } catch (error) {
    throw error
  }
}

const depositEthForOperator = async (
  web3Context,
  data,
  onTransactionHashCallback
) => {
  const { keepBondingContract, yourAddress } = web3Context
  const { operatorAddress, value } = data
  const valueInWei = web3Utils.toWei(value.toString(), "ether")

  await keepBondingContract.methods
    .deposit(operatorAddress)
    .send({ from: yourAddress, value: valueInWei })
    .on("transactionHash", onTransactionHashCallback)
}

// TODO fetch from contracts
const getBondingData = async (web3Context) => {
  const bondingData = [
    {
      operatorAddress: "0x2A489EacBf4de172B4018D2b4a405F05C400f530",
      stakeAmount: "1000",
      bondedETH: "1000",
      availableETH: "1000",
      availableETHInWei: "1000000000000000000000",
    },
  ]

  return bondingData
}

export const tbtcAuthorizationService = {
  fetchTBTCAuthorizationData,
  authorizeBondedECDSAKeepFactory,
  authorizeTBTCSystem,
  depositEthForOperator,
  getBondingData,
}
