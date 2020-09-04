import { contractService } from "./contracts.service"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
  MANAGED_GRANT_FACTORY_CONTRACT_NAME,
} from "../constants/constants"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  getBondedECDSAKeepFactoryAddress,
  getTBTCSystemAddress,
  ContractsLoaded,
} from "../contracts"
import web3Utils from "web3-utils"
import {
  getOperatorsOfAuthorizer,
  getOperatorsOfOwner,
} from "./token-staking.service"

const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()
const tBTCSystemAddress = getTBTCSystemAddress()

const fetchTBTCAuthorizationData = async (web3Context) => {
  const operatorsOfAuthorizer = await getOperatorsOfAuthorizer(
    web3Context,
    web3Context.yourAddress
  )
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

const fetchBondingData = async (web3Context) => {
  const { yourAddress } = web3Context
  const bondingData = []

  try {
    const operators = await fetchOperatorsOf(web3Context, yourAddress)
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc(web3Context)
    const createdBonds = await fetchCreatedBonds(
      web3Context,
      Array.from(operators.keys()),
      sortitionPoolAddress
    )

    const operatorBondingDataMap = new Map()
    for (let i = 0; i < createdBonds.length; i++) {
      const operatorAddress = web3Utils.toChecksumAddress(
        createdBonds[i].operator
      )
      const bondedEth = await fetchLockedBondAmount(
        web3Context,
        operatorAddress,
        createdBonds[i].holder,
        createdBonds[i].referenceID
      )

      const currentBond = operatorBondingDataMap.get(operatorAddress)
      if (currentBond) {
        operatorBondingDataMap.set(operatorAddress, add(currentBond, bondedEth))
      } else {
        operatorBondingDataMap.set(operatorAddress, bondedEth)
      }
    }

    for (const [operatorAddress, value] of operators.entries()) {
      const { isWithdrawableForOperator, managedGrantInfo } = value
      const delegatedTokens = await fetchDelegationInfo(
        web3Context,
        operatorAddress
      )
      const availableEth = await fetchAvailableAmount(
        web3Context,
        operatorAddress
      )

      const bondedEth = operatorBondingDataMap.get(operatorAddress)
        ? operatorBondingDataMap.get(operatorAddress)
        : 0

      const bonding = {
        operatorAddress,
        managedGrantAddress: managedGrantInfo.address,
        isWithdrawableForOperator,
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

const fetchSortitionPoolForTbtc = async () => {
  const { bondedEcdsaKeepFactoryContract } = await ContractsLoaded

  return await bondedEcdsaKeepFactoryContract.methods
    .getSortitionPool(tBTCSystemAddress)
    .call()
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
  let createdBonds = []
  if (!isEmptyArray(operatorAddresses)) {
    createdBonds = (
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

  return createdBonds
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
  // operatorAddress -> { managedGrantInfo: { address }, isWithdrawableForOperator: true  }
  const operators = new Map()

  // operators of grantee (yourAddress)
  const operatorsOfGrantee = await contractService.makeCall(
    web3Context,
    TOKEN_GRANT_CONTRACT_NAME,
    "getGranteeOperators",
    yourAddress
  )
  for (let i = 0; i < operatorsOfGrantee.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfGrantee[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  const managedGrantAddresses = await fetchManagedGrantAddresses(
    web3Context,
    yourAddress
  )
  for (const managedGrantAddress of managedGrantAddresses) {
    // operators of grantee (managedGrantAddress)
    const operatorsOfManagedGrant = await contractService.makeCall(
      web3Context,
      TOKEN_GRANT_CONTRACT_NAME,
      "getGranteeOperators",
      managedGrantAddress
    )
    for (const operatorOfManagedGrant of operatorsOfManagedGrant) {
      operators.set(web3Utils.toChecksumAddress(operatorOfManagedGrant), {
        managedGrantInfo: { address: managedGrantAddress },
        isWithdrawableForOperator: true,
      })
    }
  }

  // operators of authorizer
  const operatorsOfAuthorizer = await getOperatorsOfAuthorizer(
    web3Context,
    web3Context.yourAddress
  )
  for (let i = 0; i < operatorsOfAuthorizer.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfAuthorizer[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: false,
    })
  }

  // operators of owner (yourAddress as owner)
  const operatorsOfOwner = await getOperatorsOfOwner(yourAddress)

  for (let i = 0; i < operatorsOfOwner.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfOwner[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  const ownerAddress = await contractService.makeCall(
    web3Context,
    TOKEN_STAKING_CONTRACT_NAME,
    "ownerOf",
    yourAddress
  )

  if (ownerAddress !== "0x0000000000000000000000000000000000000000") {
    // yourAddress is an operator
    operators.set(web3Utils.toChecksumAddress(yourAddress), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  return operators
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
  fetchBondingData,
  fetchSortitionPoolForTbtc,
}
