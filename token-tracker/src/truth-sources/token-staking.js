/** @typedef { import("../lib/context.js").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import { ITruthSource, AddressesResolver } from "./truth-source.js"
import { Contract } from "../lib/contract-helper.js"
import { logger } from "../lib/winston.js"
import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"

import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import Web3 from "web3"
const { toBN } = Web3.utils

import TokenStakingJSON from "@keep-network/keep-core/artifacts/TokenStaking.json"

const TOKEN_STAKING_HISTORIC_STAKERS_DUMP_PATH =
  "./tmp/token-staking-stakers.json"
const TOKEN_STAKING_BALANCES_DUMP_PATH = "./tmp/token-staking-balances.json"
const TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_DUMP_PATH =
  "./tmp/token-staking-unknown_owners_contracts.json"

// TODO: Add support for the old TokenStaking contract.

export class TokenStakingTruthSource extends ITruthSource {
  constructor(
    /** @type {Context} */ context,
    /** @type {Number} */ targetBlock
  ) {
    super(context, targetBlock)
  }

  async initialize() {
    this.context.addContract(
      "TokenStaking",
      new Contract(TokenStakingJSON, this.context.web3)
    )

    this.tokenStaking = await this.context
      .getContract("TokenStaking")
      .deployed()
  }

  /**
   * Finds all historic stakers based on StakeDelegated events emitted by TokenStaking
   * contract.
   * @return {Map<Address,Address>} All historic token stake operators with their
   * owners.
   */
  async findHistoricStakeOperatorsOwners() {
    logger.info(
      `looking for StakeDelegated events emitted from ${this.tokenStaking.options.address} ` +
        `between blocks ${this.context.deploymentBlock} and ${this.targetBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.tokenStaking,
      "StakeDelegated",
      this.context.deploymentBlock,
      this.targetBlock
    )
    logger.info(`found ${events.length} stake delegated events`)

    const operatorsOwnersMap = new Map()
    events.forEach((event) => {
      operatorsOwnersMap.set(
        event.returnValues.operator,
        event.returnValues.owner
      )
    })

    dumpDataToFile(operatorsOwnersMap, TOKEN_STAKING_HISTORIC_STAKERS_DUMP_PATH)

    return operatorsOwnersMap
  }

  /**
   * Filters owners based on the rules defined for sources of truth. Ignores
   * addresses of known Keep contracts for which actual holders are resolved.
   * @param {Map<Address,Address>} operatorsOwnersMap Map of operators and their
   * owners to filter.
   * @return {Map<Address,Address>} Filtered map of operators and owners.
   */
  async filterOwners(operatorsOwnersMap) {
    logger.info(`filter owners`)

    const filteredOperatorsOwners = new Map()
    const unknownContracts = new Set()

    for (const [operator, owner] of operatorsOwnersMap) {
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

    dumpDataToFile(
      unknownContracts,
      TOKEN_STAKING_UNKNOWN_OWNERS_CONTRACTS_DUMP_PATH
    )

    return filteredOperatorsOwners
  }

  /**
   * Checks token balances for owners in the operators to owner map. It walks over
   * delegations to operators, ignored delegations that were undelegated. Combines
   * results for owners that have multiple operators.
   * @param {Map<Address,Address>} stakers Map of operators and owners to check.
   * @return {Map<Address,BN>} Token holdings at the target block.
   */
  async checkStakedValues(stakers) {
    logger.info(`check stake delegations at block ${this.targetBlock}`)

    /** @type {Map<Address,BN>} */
    const stakersBalances = new Map()

    for (const [operator, owner] of stakers) {
      const delegationInfo = await callWithRetry(
        this.tokenStaking.methods.getDelegationInfo(operator),
        undefined,
        undefined,
        this.targetBlock
      )

      const amount = toBN(delegationInfo.amount)
      const undelegatedAt = toBN(delegationInfo.undelegatedAt)

      if (undelegatedAt.gtn(0)) {
        logger.debug(
          `skipping delegation to ${operator}, undelegated at: ${undelegatedAt.toString()}`
        )
        continue
      }

      if (amount.eqn(0)) {
        logger.debug(`skipping delegation to ${operator}, amount is zero`)
        continue
      }

      if (!stakersBalances.has(owner)) {
        stakersBalances.set(owner, toBN(0))
      }
      stakersBalances.get(owner).iadd(amount)

      logger.debug(
        `owner's ${owner} total stake: ${stakersBalances.get(owner).toString()}`
      )
    }

    dumpDataToFile(stakersBalances, TOKEN_STAKING_BALANCES_DUMP_PATH)

    return stakersBalances
  }

  /**
   * Returns a map of addresses with their token holdings based on Token Staking.
   * @return {Map<Address,BN>} Token holdings at the target blocks.
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const allStakeOwners = await this.findHistoricStakeOperatorsOwners()

    const filteredStakeOwners = await this.filterOwners(allStakeOwners)

    const ownersBalances = await this.checkStakedValues(filteredStakeOwners)
    return ownersBalances
  }
}
