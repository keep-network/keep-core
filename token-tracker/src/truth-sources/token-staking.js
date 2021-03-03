/** @typedef { import("../lib/context.js").Context } Context */
/** @typedef { import("./contract-helper.js").Contract } Contract */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import { ITruthSource, AddressesResolver } from "./truth-source.js"
import { logger } from "../lib/winston.js"
import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"
import { addTokenHoldings } from "../lib/map-helper.js"

import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import Web3 from "web3"
const { toBN } = Web3.utils

const OLD_TOKEN_STAKING_HISTORIC_STAKERS_OUTPUT_PATH =
  "./tmp/token-staking-old-stakers.json"
const OLD_TOKEN_STAKING_BALANCES_OUTPUT_PATH =
  "./tmp/token-staking-old-balances.json"
const OLD_TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_OUTPUT_PATH =
  "./tmp/token-staking-old-unknown_owners_contracts.json"

const TOKEN_STAKING_HISTORIC_STAKERS_OUTPUT_PATH =
  "./tmp/token-staking-stakers.json"
const TOKEN_STAKING_BALANCES_OUTPUT_PATH = "./tmp/token-staking-balances.json"
const TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_OUTPUT_PATH =
  "./tmp/token-staking-unknown_owners_contracts.json"

export class TokenStakingTruthSource extends ITruthSource {
  /** @property {Contract} oldTokenStaking */
  /** @property {Contract} tokenStaking */

  constructor(
    /** @type {Context} */ context,
    /** @type {Number} */ targetBlock
  ) {
    super(context, targetBlock)
  }

  async initialize() {
    this.oldTokenStaking = await this.context
      .getContract("OldTokenStaking")
      .deployed()

    this.tokenStaking = await this.context
      .getContract("TokenStaking")
      .deployed()
  }

  /** @return {Map<Address,BN>} Token holdings at the target block. */
  async getHoldingsFromOldTokenStaking() {
    const allStakeOwners = await this.findHistoricStakeOperatorsOwners(
      this.oldTokenStaking,
      "Staked",
      OLD_TOKEN_STAKING_HISTORIC_STAKERS_OUTPUT_PATH
    )

    const filteredStakeOwners = await this.filterOwners(
      allStakeOwners,
      OLD_TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_OUTPUT_PATH
    )

    const ownersBalances = await this.checkStakedValues(
      this.oldTokenStaking,
      filteredStakeOwners,
      OLD_TOKEN_STAKING_BALANCES_OUTPUT_PATH
    )

    return ownersBalances
  }

  /** @return {Map<Address,BN>} Token holdings at the target block. */
  async getHoldingsFromTokenStaking() {
    const allStakeOwners = await this.findHistoricStakeOperatorsOwners(
      this.tokenStaking,
      "StakeDelegated",
      TOKEN_STAKING_HISTORIC_STAKERS_OUTPUT_PATH
    )

    const filteredStakeOwners = await this.filterOwners(
      allStakeOwners,
      TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_OUTPUT_PATH
    )

    const ownersBalances = await this.checkStakedValues(
      this.tokenStaking,
      filteredStakeOwners,
      TOKEN_STAKING_BALANCES_OUTPUT_PATH
    )

    return ownersBalances
  }

  /**
   * Finds all historic stakers based on staked events emitted by TokenStaking
   * contract.
   * @param {Contract} stakingContract Staking contract to verify.
   * @param {String} stakedEventName Name of the event emitted on stake.
   * @param {String} stakersOutputPath Path to a file where result should be stored.
   * @return {Map<Address,Address>} All historic token stake operators with their
   * owners.
   */
  async findHistoricStakeOperatorsOwners(
    stakingContract,
    stakedEventName,
    stakersOutputPath
  ) {
    logger.info(
      `looking for ${stakedEventName} events emitted from ${stakingContract.options.address} ` +
        `between blocks ${this.context.deploymentBlock} and ${this.targetBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      stakingContract,
      stakedEventName,
      this.context.deploymentBlock,
      this.targetBlock
    )
    logger.info(`found ${events.length} ${stakedEventName} events`)

    const operatorsOwnersMap = new Map()
    events.forEach((event) => {
      // Events emitted on stake differs between the OldTokenStaking and TokenStaking:
      // OldTokenStaking: `Staked(address indexed from, uint256 value)` (where from is an operator)
      // TokenStaking: `StakeDelegated(address indexed owner, address indexed operator)`
      operatorsOwnersMap.set(
        event.returnValues.operator || event.returnValues.from,
        event.returnValues.owner
      )
    })

    dumpDataToFile(operatorsOwnersMap, stakersOutputPath)

    return operatorsOwnersMap
  }

  /**
   * Filters owners based on the rules defined for sources of truth. Ignores
   * addresses of known Keep contracts for which actual holders are resolved.
   * @param {Map<Address,Address>} operatorsOwnersMap Map of operators and their
   * owners to filter.
   * @param {String} unknownOwnersOutputPath Path to a file where result should be stored.
   * @return {Map<Address,Address>} Filtered map of operators and owners.
   */
  async filterOwners(operatorsOwnersMap, unknownOwnersOutputPath) {
    logger.info(`filter owners`)

    const filteredOperatorsOwners = new Map()
    const unknownContracts = new Set()

    for (let [operator, owner] of operatorsOwnersMap) {
      if (!owner) {
        logger.debug(
          `owner not provided for operator ${operator}; ` +
            `fetching it now from the OldTokenStaking contract`
        )

        owner = await callWithRetry(
          this.oldTokenStaking.methods.ownerOf(operator),
          undefined,
          undefined,
          this.targetBlock
        )
      }

      const [
        isIgnored,
        addressType,
      ] = await this.addressesResolver.isIgnoredAddress(owner)

      logger.debug(
        `operator's [${operator}] owner [${owner}] is ${addressType}`
      )

      if (isIgnored) {
        continue // skip this operator
      }

      // If the address doesn't match any of ignored contracts store it for
      // reference if we want to double check them.
      if (addressType === AddressesResolver.UNKNOWN_CONTRACT)
        unknownContracts.add(owner)

      filteredOperatorsOwners.set(operator, owner)
    }

    dumpDataToFile(unknownContracts, unknownOwnersOutputPath)

    return filteredOperatorsOwners
  }

  /**
   * Checks token balances for owners in the operators to owner map. It walks over
   * delegations to operators, ignored delegations that were undelegated. Combines
   * results for owners that have multiple operators.
   * @param {Contract} stakingContract Staking contract to verify.
   * @param {Map<Address,Address>} stakers Map of operators and owners to check.
   * @param {String} balancesOutputPath Path to a file where result should be stored.
   * @return {Map<Address,BN>} Token holdings at the target block.
   */
  async checkStakedValues(stakingContract, stakers, balancesOutputPath) {
    logger.info(`check stake delegations at block ${this.targetBlock}`)

    /** @type {Map<Address,BN>} */
    const stakersBalances = new Map()

    for (const [operator, owner] of stakers) {
      const delegationInfo = await callWithRetry(
        stakingContract.methods.getDelegationInfo(operator),
        undefined,
        undefined,
        this.targetBlock
      )

      const amount = toBN(delegationInfo.amount)

      if (amount.eqn(0)) {
        logger.debug(`skipping delegation to ${operator}, amount is zero`)
        continue
      }

      if (stakersBalances.has(owner)) {
        stakersBalances.get(owner).iadd(amount)
      } else {
        stakersBalances.set(owner, amount)
      }

      logger.debug(
        `owner ${owner} staked ${amount} to operator ${operator}; owner's total stake: ${stakersBalances
          .get(owner)
          .toString()}`
      )
    }

    dumpDataToFile(stakersBalances, balancesOutputPath)

    return stakersBalances
  }

  /**
   * Returns a map of addresses with their token holdings based on Token Staking.
   * @return {Map<Address,BN>} Token holdings at the target blocks.
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const oldTokenStakingBalances = await this.getHoldingsFromOldTokenStaking()
    const tokenStakingBalances = await this.getHoldingsFromTokenStaking()

    const result = addTokenHoldings(
      oldTokenStakingBalances,
      tokenStakingBalances
    )

    return result
  }
}
