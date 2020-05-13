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
} from "../contracts"
import web3Utils from "web3-utils"

const tBTCSystemAddress = getTBTCSystemAddress()
const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()

const fetchTBTCAuthorizationData = async (web3Context) => {
  const { yourAddress } = web3Context

  const stakedEvents = await fetchStakedEvents(web3Context)
  const visitedOperators = {}
  const authorizerOperators = []

  // TODO: remove console logs
  console.log(
    "getBondedECDSAKeepFactoryAddress: ",
    getBondedECDSAKeepFactoryAddress()
  )
  console.log("getTBTCSystemAddress: ", getTBTCSystemAddress())

  // Fetch all authorizer operators
  for (let i = 0; i < stakedEvents.length; i++) {
    const {
      returnValues: { from: operatorAddress },
    } = stakedEvents[i]

    if (visitedOperators.hasOwnProperty(operatorAddress)) continue

    visitedOperators[operatorAddress] = operatorAddress
    const authorizerOfOperator = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "authorizerOf",
      operatorAddress
    )

    if (isSameEthAddress(authorizerOfOperator, yourAddress)) {
      const delegatedTokens = await fetchDelegationInfo(
        web3Context,
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

const getBondingData = async (web3Context) => {
  const { yourAddress } = web3Context

  const bondingData = []

  try {
    const operators = await fetchOperatorsOf(web3Context, yourAddress)
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc(web3Context)
    const createdBondsEvents = await fetchCreatedBondsEvents(
      web3Context,
      operators,
      sortitionPoolAddress
    )

    const operatorBondingDataMap = new Map()
    if (createdBondsEvents.length !== 0) {
      for (let i = 0; i < createdBondsEvents.length; i++) {
        const {
          returnValues: { operator, holder, referenceID },
        } = createdBondsEvents[i]

        let bondHolder = holder
        let bondReferenceId = referenceID

        const reassignedEvents = await fetchBondReassignedEvents(
          web3Context,
          bondHolder,
          bondReferenceId
        )
        if (reassignedEvents.length > 0) {
          // TODO: need to test reasssignment
          const latestReassignedEvent =
            reassignedEvents[reassignedEvents.length() - 1]
          const {
            returnValues: { holder, referenceID },
          } = latestReassignedEvent

          bondHolder = holder
          bondReferenceId = referenceID
        }

        const bondedEth = await fetchLockedBondAmount(
          web3Context,
          operator,
          bondHolder,
          bondReferenceId
        )

        operatorBondingDataMap.set(operator, bondedEth)
      }
    }

    for (let i = 0; i < operators.length; i++) {
      const delegatedTokens = await fetchDelegationInfo(
        web3Context,
        operators[i]
      )
      const availableEth = await fetchAvailableAmount(
        web3Context,
        operators[i],
        bondedECDSAKeepFactoryAddress,
        sortitionPoolAddress
      )

      let bondedEth = 0
      if (operatorBondingDataMap.get(operators[i])) {
        bondedEth = operatorBondingDataMap.get(operators[i])
      }

      const bonding = {
        operatorAddress: operators[i],
        stakeAmount: delegatedTokens.amount,
        bondedETH: web3Utils.fromWei(bondedEth.toString(), "ether"),
        availableETH: web3Utils.fromWei(availableEth.toString(), "ether"),
        availableETHInWei: availableEth,
      }

      bondingData.push(bonding)
    }
  } catch (error) {
    // return error / false?
    console.error("failed to fetch bonds for tBTC", error)
  }

  return bondingData
}

const fetchStakedEvents = async (web3Context) => {
  return contractService.getPastEvents(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "Staked",
    { fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME] }
  )
}

const fetchSortitionPoolForTbtc = async (web3Context) => {
  return contractService.makeCall(
    web3Context,
    BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
    "getSortitionPool",
    tBTCSystemAddress
  )
}

const fetchDelegationInfo = async (web3Context, operatorAddress) => {
  return contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "getDelegationInfo",
    operatorAddress
  )
}

const fetchCreatedBondsEvents = async (
  web3Context,
  operatorAddresses,
  sortitionPoolAddress
) => {
  return contractService.getPastEvents(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "BondCreated",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[KEEP_BONDING_CONTRACT_NAME],
      filter: {
        operator: operatorAddresses,
        sortitionPool: sortitionPoolAddress,
      },
    }
  )
}

const fetchBondReassignedEvents = async (
  web3Context,
  operatorAddresses,
  referenceId
) => {
  return contractService.getPastEvents(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "BondReassigned",
    {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[KEEP_BONDING_CONTRACT_NAME],
      filter: { operator: operatorAddresses, referenceID: referenceId },
    }
  )
}

const fetchOperatorsOf = async (web3Context, yourAddress) => {
  const ownerOperators = await contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "operatorsOf",
    yourAddress
  )

  if (ownerOperators.length === 0) {
    ownerOperators[0] = yourAddress
  }

  return ownerOperators
}

// aka lockedBonds
const fetchLockedBondAmount = async (
  web3Context,
  operator,
  holder,
  referenceID
) => {
  return contractService.makeCall(
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
  bondedECDSAKeepFactoryAddress,
  authorizedSortitionPool
) => {
  return contractService.makeCall(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "availableUnbondedValue",
    operator,
    bondedECDSAKeepFactoryAddress,
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
