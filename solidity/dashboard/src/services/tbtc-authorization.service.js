import { contractService } from "./contracts.service"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
} from "../constants/constants"
import {
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
} from "../constants/constants"
import { isSameEthAddress } from "../utils/general.utils"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  getBondedECDSAKeepFactoryAddress,
  getTBTCSystemAddress,
} from "../contracts"
import web3Utils from "web3-utils"

const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()
const tBTCSystemAddress = getTBTCSystemAddress()

const fetchTBTCAuthorizationData = async (web3Context) => {
  const operatorsOfAuthorizer = await fetchOperatorsOfAuthorizer(web3Context)
  const tbtcAuthorizatioData = []

  for (let i = 0; i < operatorsOfAuthorizer.length; i++) {
    const delegatedTokens = await fetchDelegationInfo(
      web3Context,
      operatorsOfAuthorizer[i]
    )

    const isBondedECDSAKeepFactoryAuthorized = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "isAuthorizedForOperator",
      operatorsOfAuthorizer[i],
      bondedECDSAKeepFactoryAddress
    )

    const isTBTCSystemAuthorized = await isTbtcSystemAuthorized(
      web3Context,
      operatorsOfAuthorizer[i]
    )

    const authorizerOperator = {
      operatorAddress: operatorsOfAuthorizer[i],
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

    tbtcAuthorizatioData.push(authorizerOperator)
  }

  return tbtcAuthorizatioData
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

const withdrawAllEthForOperator = async (
  web3Context,
  data,
  onTransactionHashCallback
) => {
  const { keepBondingContract, yourAddress } = web3Context
  const { operatorAddress, availableETH } = data
  const availableInWei = web3Utils.toWei(availableETH.toString(), "ether")

  await keepBondingContract.methods
    .withdraw(availableInWei, operatorAddress)
    .send({ from: yourAddress })
    .on("transactionHash", onTransactionHashCallback)
}

const fetchOperatorsOfAuthorizer = async (web3Context) => {
  const { yourAddress } = web3Context
  const stakedOperatorAddresses = await fetchStakedOperators(web3Context)
  const visitedOperators = {}
  const operatorsOfAuthorizer = []

  for (let i = 0; i < stakedOperatorAddresses.length; i++) {
    if (visitedOperators.hasOwnProperty(stakedOperatorAddresses[i])) continue
    visitedOperators[stakedOperatorAddresses[i]] = stakedOperatorAddresses[i]

    const authorizerOfOperator = await fetchAuthorizerOfOperator(
      web3Context,
      stakedOperatorAddresses[i]
    )

    // Accept only operators of an authorizer
    if (isSameEthAddress(authorizerOfOperator, yourAddress)) {
      operatorsOfAuthorizer.push(stakedOperatorAddresses[i])
    }
  }

  return operatorsOfAuthorizer
}

const fetchBondingData = async (web3Context) => {
  const { yourAddress } = web3Context
  const bondingData = []

  try {
    const operators = await fetchOperatorsOf(web3Context, yourAddress)
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc(web3Context)
    const createdBonds = await fetchCreatedBonds(
      web3Context,
      operators.map((_) => _[0]),
      sortitionPoolAddress
    )

    const operatorBondingDataMap = new Map()
    for (let i = 0; i < createdBonds.length; i++) {
      const bondedEth = await fetchLockedBondAmount(
        web3Context,
        createdBonds[i].operator,
        createdBonds[i].holder,
        createdBonds[i].referenceID
      )

      operatorBondingDataMap.set(
        web3Utils.toChecksumAddress(createdBonds[i].operator),
        bondedEth
      )
    }

    for (let i = 0; i < operators.length; i++) {
      const delegatedTokens = await fetchDelegationInfo(
        web3Context,
        operators[i][0]
      )
      const availableEth = await fetchAvailableAmount(
        web3Context,
        operators[i][0]
      )

      const bondedEth = operatorBondingDataMap.get(
        web3Utils.toChecksumAddress(operators[i][0])
      )
        ? operatorBondingDataMap.get(
            web3Utils.toChecksumAddress(operators[i][0])
          )
        : 0

      const bonding = {
        operatorAddress: operators[i][0],
        isWithdrawable: operators[i][1],
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

const fetchStakedOperators = async (web3Context) => {
  return (
    await contractService.getPastEvents(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "Staked",
      { fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER[TOKEN_STAKING_CONTRACT_NAME] }
    )
  ).map((_) => _.returnValues.from)
}

const fetchAuthorizerOfOperator = async (web3Context, operatorAddress) => {
  return contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "authorizerOf",
    operatorAddress
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

const fetchCreatedBonds = async (
  web3Context,
  operatorAddresses,
  sortitionPoolAddress
) => {
  return (
    await contractService.getPastEvents(
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
  ).map((_) => {
    return {
      operator: _.returnValues.operator,
      holder: _.returnValues.holder,
      referenceID: _.returnValues.referenceID,
    }
  })
}

const fetchManagedGrantAddresses = async (web3Context, lookupAddress) => {
  return (
    await contractService.getPastEvents(
      web3Context,
      MANAGED_GRANT_FACTORY_CONTRACT_NAME,
      "ManagedGrantCreated",
      {
        fromBlock:
          CONTRACT_DEPLOY_BLOCK_NUMBER[MANAGED_GRANT_FACTORY_CONTRACT_NAME],
        filter: { grantee: lookupAddress },
      }
    )
  ).map((_) => _.returnValues.grantAddress)
}

const fetchOperatorsOf = async (web3Context, yourAddress) => {
  const operators = new Map()

  // operators of grantee (yourAddress)
  const operatorsOfGrantee = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "getGranteeOperators",
    yourAddress
  )
  for (let i = 0; i < operatorsOfGrantee.length; i++) {
    operators.set(operatorsOfGrantee[i], false)
  }

  const managedGrantAddresses = await fetchManagedGrantAddresses(
    web3Context,
    yourAddress
  )
  for (let i = 0; i < managedGrantAddresses.length; ++i) {
    const managedGrantAddress = managedGrantAddresses[i]
    // operators of grantee (managedGrantAddress)
    const operatorsOfManagedGrant = await contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "getGranteeOperators",
      managedGrantAddress
    )
    for (let i = 0; i < operatorsOfManagedGrant.length; i++) {
      operators.set(operatorsOfManagedGrant[i], false)
    }
  }

  // operators of authorizer
  const operatorsOfAuthorizer = await fetchOperatorsOfAuthorizer(web3Context)
  for (let i = 0; i < operatorsOfAuthorizer.length; i++) {
    operators.set(operatorsOfAuthorizer[i], false)
  }

  // operators of owner
  const operatorsOfOwner = await contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "operatorsOf",
    yourAddress // as owner
  )
  for (let i = 0; i < operatorsOfOwner.length; i++) {
    operators.set(operatorsOfOwner[i], true)
  }

  if (operators.length === 0) {
    const ownerAddress = await contractService.makeCall(
      web3Context,
      TOKEN_STAKING_CONTRACT_NAME,
      "ownerOf",
      yourAddress
    )

    if (ownerAddress !== "0x0000000000000000000000000000000000000000") {
      // yourAddress is an operator
      operators.set(yourAddress, true)
    }
  }

  return [...operators]
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
const fetchAvailableAmount = async (web3Context, operator) => {
  return contractService.makeCall(
    web3Context,
    KEEP_BONDING_CONTRACT_NAME,
    "unbondedValue",
    operator
  )
}

export const tbtcAuthorizationService = {
  fetchTBTCAuthorizationData,
  authorizeBondedECDSAKeepFactory,
  authorizeTBTCSystem,
  fetchBondingData,
  depositEthForOperator,
  withdrawAllEthForOperator,
}
