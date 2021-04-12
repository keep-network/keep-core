import {
  KEEP_BONDING_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  STAKING_PORT_BACKER_CONTRACT_NAME,
} from "../constants/constants"
import { add } from "../utils/arithmetics.utils"
import { isEmptyArray } from "../utils/array.utils"
import {
  getContractDeploymentBlockNumber,
  getBondedECDSAKeepFactoryAddress,
  getTBTCSystemAddress,
  ContractsLoaded,
} from "../contracts"
import web3Utils from "web3-utils"
import {
  getOperatorsOfAuthorizer,
  getOperatorsOfOwner,
} from "./token-staking.service"
import { tokenGrantsService } from "./token-grants.service"
import { ZERO_ADDRESS } from "../utils/ethereum.utils"

const bondedECDSAKeepFactoryAddress = getBondedECDSAKeepFactoryAddress()
const tBTCSystemAddress = getTBTCSystemAddress()

const fetchTBTCAuthorizationData = async (address) => {
  if (!address) {
    return []
  }

  const { stakingContract } = await ContractsLoaded
  const operatorsOfAuthorizer = await getOperatorsOfAuthorizer(address)
  const tbtcAuthorizatioData = []

  for (let i = 0; i < operatorsOfAuthorizer.length; i++) {
    const delegatedTokens = await fetchDelegationInfo(operatorsOfAuthorizer[i])

    const isBondedECDSAKeepFactoryAuthorized = await stakingContract.methods
      .isAuthorizedForOperator(
        operatorsOfAuthorizer[i],
        bondedECDSAKeepFactoryAddress
      )
      .call()

    const isTBTCSystemAuthorized = await isTbtcSystemAuthorized(
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

const isTbtcSystemAuthorized = async (operatorAddress) => {
  const {
    bondedEcdsaKeepFactoryContract,
    keepBondingContract,
  } = await ContractsLoaded
  try {
    const sortitionPoolAddress = await bondedEcdsaKeepFactoryContract.methods
      .getSortitionPool(tBTCSystemAddress)
      .call()

    return await keepBondingContract.methods
      .hasSecondaryAuthorization(operatorAddress, sortitionPoolAddress)
      .call()
  } catch {
    return false
  }
}

const fetchBondingData = async (address) => {
  const bondingData = []
  if (!address) {
    return bondingData
  }

  try {
    const operators = await fetchOperatorsOf(address)
    const sortitionPoolAddress = await fetchSortitionPoolForTbtc()
    const createdBonds = await fetchCreatedBonds(
      Array.from(operators.keys()),
      sortitionPoolAddress
    )

    const operatorBondingDataMap = new Map()
    for (let i = 0; i < createdBonds.length; i++) {
      const operatorAddress = web3Utils.toChecksumAddress(
        createdBonds[i].operator
      )
      const bondedEth = await fetchLockedBondAmount(
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
      const delegatedTokens = await fetchDelegationInfo(operatorAddress)
      const availableEth = await fetchAvailableAmount(operatorAddress)

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

const fetchDelegationInfo = async (operatorAddress) => {
  const { stakingContract } = await ContractsLoaded
  return await stakingContract.methods.getDelegationInfo(operatorAddress).call()
}

const fetchCreatedBonds = async (operatorAddresses, sortitionPoolAddress) => {
  const { keepBondingContract } = await ContractsLoaded
  let createdBonds = []
  if (!isEmptyArray(operatorAddresses)) {
    createdBonds = (
      await keepBondingContract.getPastEvents("BondCreated", {
        fromBlock: await getContractDeploymentBlockNumber(
          KEEP_BONDING_CONTRACT_NAME
        ),
        filter: {
          operator: operatorAddresses,
          sortitionPool: sortitionPoolAddress,
        },
      })
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

const fetchOperatorsOf = async (address) => {
  const { grantContract, stakingContract } = await ContractsLoaded
  /**
   * Operator address to details.
   * @type {Map<string, { managedGrantInfo: { address: string }, isWithdrawableForOperator: boolean }>}
   */
  const operators = new Map()

  // operators of authorizer
  const operatorsOfAuthorizer = await getOperatorsOfAuthorizer(address)
  for (let i = 0; i < operatorsOfAuthorizer.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfAuthorizer[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: false,
    })
  }

  // operators of grantee (yourAddress)
  const operatorsOfGrantee = await getGranteeOperators(address)

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
  const operatorsOfOwner = await getOperatorsOfOwner(address)

  for (let i = 0; i < operatorsOfOwner.length; i++) {
    operators.set(web3Utils.toChecksumAddress(operatorsOfOwner[i]), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  const copiedOperatorsFromLiquidTokens = await getCopiedOperatorsFromLiquidTokens(
    address,
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

  const ownerAddress = await stakingContract.methods.ownerOf(address).call()

  if (ownerAddress !== ZERO_ADDRESS) {
    // yourAddress is an operator
    operators.set(web3Utils.toChecksumAddress(address), {
      managedGrantInfo: {},
      isWithdrawableForOperator: true,
    })
  }

  return operators
}

// aka lockedBonds
const fetchLockedBondAmount = async (operator, holder, referenceID) => {
  const { keepBondingContract } = await ContractsLoaded
  return await keepBondingContract.methods
    .bondAmount(operator, holder, referenceID)
    .call()
}

// aka unbondedValue
const fetchAvailableAmount = async (operator) => {
  const { keepBondingContract } = await ContractsLoaded
  return await keepBondingContract.methods.unbondedValue(operator).call()
}

const getGranteeOperators = async (address) => {
  const { grantContract } = await ContractsLoaded

  // Fetch all grantee operators. These are not all grantee operators,
  // since `TokenGrant` contract does not know about escrow redelegation.
  const operatorsOfGrantee = await grantContract.methods
    .getGranteeOperators(address)
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
      fromBlock: await getContractDeploymentBlockNumber(
        TOKEN_STAKING_ESCROW_CONTRACT_NAME
      ),
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
        fromBlock: await getContractDeploymentBlockNumber(
          TOKEN_STAKING_ESCROW_CONTRACT_NAME
        ),
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
          fromBlock: await getContractDeploymentBlockNumber(
            TOKEN_STAKING_CONTRACT_NAME
          ),
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
      fromBlock: await getContractDeploymentBlockNumber(
        STAKING_PORT_BACKER_CONTRACT_NAME
      ),
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
