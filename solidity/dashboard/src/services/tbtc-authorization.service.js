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

  const stakedEvents = await fetchStakedEvents(web3Context)
  const visitedOperators = {}
  const authorizerOperators = []

  //TODO: remove console logs
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
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc(web3Context)

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
  const { yourAddress } = web3Context

  const bondingData = []
  const bondedEth = 0

  try {
    const operators = await fetchOperatorsOf(web3Context, yourAddress)
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc(web3Context)
    const createdBondsEvents = await fetchCreatedBondsEvents(web3Context, operators, sortitionPoolAddress)
    console.log("createdBondsEvents: ", createdBondsEvents)
    
    // create a map
    // operator -> bondedEth
    const operatorBondedEthMap = new Map()
    if (createdBondsEvents.length != 0) {
      for (let i = 0; i < createdBondsEvents.length; i++) {

        // const {
        //   returnValues: { from: operator },
        // } = createdBondsEvents[i]

        // holder = createBondsEvents[i].holder
        // referenceID = createBondsEvents[i].referenceID

        // const reassignedEvents = fetchBondReassignedEvents(web3Context, operator, referenceId) 
        // if reassignedEventsForOperator > 0 
        //  latestReassignedEvent = reassignedEvents[reassignedEvents.length() - 1]
        //  holder = latestReassignedEvent[i].holder
        //  referenceID = latestReassignedEvent[i].referenceID


        // bondedEth = fetchLockedBondAmount(operator, holder, referenceID)

        //operatorBondedEthMap.set(operator, bondedEth)
      }
    } 

    for (let i = 0; i < operators.length; i++) {
      
      console.log("operator for bonding TBTC: ", operators[i])

      const delegatedTokens = await contractService.makeCall(
        web3Context,
        TOKEN_STAKING_CONTRACT_NAME,
        "getDelegationInfo",
        yourAddress,
      )

      // if (operatorBondedEthMap.get(operator) != undefined)
      //  bondedEth = operatorBondedEthMap.get(operator)

      const bonding = {
        operatorAddress: operators[i],
        stakeAmount: delegatedTokens.amount,
        bondedETH: bondedEth,
        availableETH: {},
        availableETHInWei: "1000000000000000000000",
      }

      bondingData.push(bonding)
    }

  } catch(error) {
    // return error / false?
    console.error("failed to fetch bonds for tBTC", error)
  }

  return bondingData
}

const fetchStakedEvents = async (web3Context) => {
  return await contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "Staked",
    { fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME] }
  )
}

const fetchSortitionPoolForTbtc = async (web3Context) => {
  return await contractService.makeCall(
    web3Context,
    BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
    "getSortitionPool",
    tBTCSystemAddress,
  )
}

const fetchCreatedBondsEvents = async (
  web3Context, 
  operatorAddresses, 
  sortitionPoolAddress
  ) => {

  return await contractService.getPastEvents(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "BondCreated",
    { 
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[KEEP_BONDING_CONTRACT_NAME],
      filter: {operator: operatorAddresses, sortitionPool: sortitionPoolAddress},
    }
  )
}

const fetchBondReassignedEvents = async (
  web3Context, 
  operatorAddresses, 
  referenceId
  ) => {

  return await contractService.getPastEvents(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "BondReassigned",
    { 
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[KEEP_BONDING_CONTRACT_NAME],
      filter: {operator: operatorAddresses, referenceID: referenceId},
    }
  )
}

const fetchOperatorsOf = async (
  web3Context,
  yourAddress
  ) => {

  return await contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "operatorsOf",
    yourAddress,
  )
}

// aka lockedBonds
const fetchLockedBondAmount = async (
  web3Context,
  operator,
  holder,
  referenceID
  ) => {

  return await contractService.makeCall(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "bondAmount",
    operator,
    holder,
    referenceID
  )
}

// aka unbondedValue
const fetchAvailableAmount = async (
  web3Context,
  operator,
  authorizedSortitionPool
  ) => {
  return await contractService.makeCall(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "availableUnbondedValue",
    operator,
    getBondedECDSAKeepFactoryAddress(),
    authorizedSortitionPool
  )
}


export const tbtcAuthorizationService = {
  fetchTBTCAuthorizationData,
  authorizeBondedECDSAKeepFactory,
  authorizeTBTCSystem,
  depositEthForOperator,
  getBondingData,
}
