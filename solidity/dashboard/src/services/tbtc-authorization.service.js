import { contractService } from "./contracts.service"
import {
  TOKEN_STAKING_CONTRACT_NAME,
  BONDED_ECDSA_KEEP_FACTORY_CONTRACT_NAME,
  KEEP_BONDING_CONTRACT_NAME,
} from "../constants/constants"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import {
  CONTRACT_DEPLOY_BLOCK_NUMBER,
  getBondedECDSAKeepFactoryAddress,
  getTBTCSystemAddress,
  ContractsLoaded,
  Web3Loaded,
} from "../contracts"
import web3Utils from "web3-utils"
import {
  getOperatorsOfAuthorizer,
  getOperatorsOfOwner,
} from "./token-staking.service"
import { tokenGrantsService } from "./token-grants.service"

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

const fetchOperatorsOf = async (web3Context, yourAddress) => {
  const {
    eth: { defaultAccount },
  } = await Web3Loaded
  const { grantContract } = await ContractsLoaded
  /**
   * Operator address to details.
   * @type {Map<string, { managedGrantInfo: { address: string }, isWithdrawableForOperator: boolean }>}
   */
  const operators = new Map()

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

  // operators of grantee (yourAddress)
  const operatorsOfGrantee = await getGranteeOperators()

  for (let i = 0; i < operatorsOfGrantee.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfGrantee[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  const managedGrants = await tokenGrantsService.fetchManagedGrants()

  for (const managedGrant of managedGrants) {
    const managedGrantAddress =
      managedGrant.managedGrantContractInstance.options.address
    // operators of grantee (managedGrantAddress)
    const operatorsOfManagedGrant = await grantContract.methods
      .getGranteeOperators(managedGrantAddress)
      .call()
    const allOperators = await getAllGranteeOperators(operatorsOfManagedGrant)
    for (const operatorOfManagedGrant of allOperators) {
      operators.set(web3Utils.toChecksumAddress(operatorOfManagedGrant), {
        managedGrantInfo: { address: managedGrantAddress },
        isWithdrawableForOperator: true,
      })
    }
  }

  // operators of owner (yourAddress as owner)
  const operatorsOfOwner = await getOperatorsOfOwner(yourAddress)

  for (let i = 0; i < operatorsOfOwner.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfOwner[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  const copiedOperatorsFromLiquidTokens = await getCopiedOperatorsFromLiquidTokens(
    defaultAccount,
    Array.from(operators.keys())
  )
  for (let i = 0; i < copiedOperatorsFromLiquidTokens.length; i++) {
    operators.set(
      web3Utils.toChecksumAddress(copiedOperatorsFromLiquidTokens[i]),
      {
        managedGrantInfo: {},
        // From the `TokenStaking` contract's perspective,
        // the `StakingPortBacker` contract is an owner of the copied delegation,
        // so the "real" owner cannot withdraw bond.
        isWithdrawableForOperator: false,
      }
    )
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

const getGranteeOperators = async () => {
  const web3 = await Web3Loaded
  const { defaultAccount } = web3.eth
  const { grantContract } = await ContractsLoaded

  // Fetch all grantee operators. These are not all grantee operators,
  // since `TokenGrant` contract does not know about escrow redelegation.
  const operatorsOfGrantee = await grantContract.methods
    .getGranteeOperators(defaultAccount)
    .call()

  return await getAllGranteeOperators(operatorsOfGrantee)
}

const getAllGranteeOperators = async (operatorsOfGrantee) => {
  const { stakingContract, tokenStakingEscrow } = await ContractsLoaded

  if (isEmptyArray(operatorsOfGrantee)) {
    return []
  }

  // We need to take into account that the delegation from a grant can be redelegated to a new operator.
  const grantIdsToScan = (
    await tokenStakingEscrow.getPastEvents("DepositRedelegated", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
      filter: {
        previousOperator: operatorsOfGrantee,
      },
    })
  ).map((_) => _.returnValues.grantId)

  let activeOperators = []
  const newOperators = []
  const previousOperators = []
  const redelagations = isEmptyArray(grantIdsToScan)
    ? []
    : await tokenStakingEscrow.getPastEvents("DepositRedelegated", {
        fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.tokenStakingEscrow,
        filter: {
          grantId: grantIdsToScan,
        },
      })

  redelagations.forEach((redelegation) => {
    newOperators.push(redelegation.returnValues.newOperator)
    previousOperators.push(redelegation.returnValues.previousOperator)
  })

  newOperators.forEach((operator) => {
    const indexOf = previousOperators.indexOf(operator)
    if (indexOf > -1) {
      previousOperators.splice(indexOf, 1)
    } else {
      activeOperators.push(operator)
    }
  })

  // Filter out obsolete operators from `TokenGrant::getGranteeOperators`.
  const activeOperatorsFromTokenGrant = operatorsOfGrantee.filter(
    (operator) => !previousOperators.includes(operator)
  )
  activeOperators = activeOperators.concat(activeOperatorsFromTokenGrant)

  return isEmptyArray(activeOperators)
    ? activeOperators
    : // Scan `OperatorStaked` events to make sure we only return operators, from a "new" `TokenStaking`,
      // since `TokenGrant` stores the old grantee-operator relationship.
      (
        await stakingContract.getPastEvents("OperatorStaked", {
          fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingContract,
          filter: { operator: activeOperators },
        })
      ).map((_) => _.returnValues.operator)
}

const getCopiedOperatorsFromLiquidTokens = async (
  ownerOrGrantee,
  operatorsToFilterOut
) => {
  const { stakingPortBackerContract } = await ContractsLoaded

  const operatorsToCheck = (
    await stakingPortBackerContract.getPastEvents("StakeCopied", {
      fromBlock: CONTRACT_DEPLOY_BLOCK_NUMBER.stakingPortBackerContract,
      filter: { owner: ownerOrGrantee },
    })
  )
    .map((_) => _.returnValues.operator)
    .filter(
      (operatorAddress) => !operatorsToFilterOut.includes(operatorAddress)
    )

  // No operators.
  if (isEmptyArray(operatorsToCheck)) {
    return []
  }

  // We only want operators for whom the delegation has not been paid back.
  const operatorsOfPortBacker = await getOperatorsOfOwner(
    stakingPortBackerContract.options.address,
    operatorsToCheck
  )

  return operatorsOfPortBacker
}

export const tbtcAuthorizationService = {
  fetchTBTCAuthorizationData,
  fetchBondingData,
  fetchSortitionPoolForTbtc,
}
