/** @typedef { import("../lib/context.js").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import { ITruthSource } from "./index.js"
import { Contract } from "../lib/contract-helper.js"

import { getPastEvents, callWithRetry } from "../lib/ethereum-helper.js"
import { writeFileSync, readFileSync } from "fs"
import BN from "bn.js"

import TokenStakingJson from "@keep-network/keep-core/artifacts/TokenStaking.json"
import { mapToObject } from "../lib/map-helper.js"

// import { toBN } from "web3-utils" FIXME: workaround
import web3Utils from "web3-utils"
const { toBN } = web3Utils

const TOKEN_STAKING_HISTORIC_STAKERS_DUMP_PATH =
  "./tmp/token-staking-stakers.json"
const TOKEN_STAKING_BALANCES_DUMP_PATH = "./tmp/token-staking-balances.json"

/**
 * TODO: Write docs
 * Short description.
 * @typedef {ITruthSource} TokenStakingTruthSource
 */

export class TokenStakingTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} finalBlock
   */
  constructor(context, finalBlock) {
    super(context, finalBlock)
  }

  async initialize() {
    const TokenStaking = new Contract(TokenStakingJson, this.context.web3)

    this.context.addContract("TokenStaking", TokenStaking)

    this.tokenStaking = await this.context.contracts.TokenStaking.deployed()
  }

  /**
   * @returns {Map<Address,Address>} All historic token stake operators with owners.
   */
  async findHistoricStakeOperatorsOwners() {
    console.info(
      `looking for StakeDelegated events emitted from ${this.tokenStaking.options.address} ` +
        `between blocks ${this.context.contracts.deploymentBlock} and ${this.finalBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.tokenStaking,
      "StakeDelegated",
      this.context.contracts.deploymentBlock,
      this.finalBlock
    )
    console.info(`found ${events.length} stake delegated events`)

    const operatorsOwnersMap = new Map()
    events.forEach((event) => {
      operatorsOwnersMap.set(
        event.returnValues.operator,
        event.returnValues.owner
      )
    })

    console.info(
      `dump all historic stakes to a file: ${TOKEN_STAKING_HISTORIC_STAKERS_DUMP_PATH}`
    )
    writeFileSync(
      TOKEN_STAKING_HISTORIC_STAKERS_DUMP_PATH,
      JSON.stringify(mapToObject(operatorsOwnersMap), null, 2)
    )

    return operatorsOwnersMap
  }

  /**
   * @param {Map<Address,Address>} stakes Stake owners to filter.
   * @returns {Map<Address,Address>} .
   */
  async filterEOA(stakes) {
    // FIXME: If the owner is TokenGrantStake or TokenStakingEscrow or StakingPortBacker contract, ignore the delegation as this amount has been already included in TokenGrant check
    const filteredEOAs = new Map()
    for (const [operator, owner] of stakes) {
      const code = await this.context.web3.eth.getCode(owner)

      if (code == "0x") filteredEOAs.set(operator, owner)
    }

    return filteredEOAs
  }

  /**
   * @param {Map<Address,Address>} stakers Token holders to check.
   * @returns {Map<Address,BN} Token holdings at the final blocks.
   */
  async checkStakedValues(stakers) {
    console.info(`check stake delegations at block ${this.finalBlock}`)

    /** @type {Map<Address,BN>} */
    const stakersBalances = new Map()

    for (const [operator, owner] of stakers) {
      const delegationInfo = await callWithRetry(
        this.tokenStaking.methods.getDelegationInfo(operator),
        undefined,
        undefined,
        this.finalBlock
      )

      const amount = toBN(delegationInfo.amount)
      const undelegatedAt = toBN(delegationInfo.undelegatedAt)

      if (undelegatedAt.gtn(0)) {
        console.debug(
          `skipping delegation to ${operator}, undelegated at: ${undelegatedAt.toString()}`
        )
        continue
      }

      if (amount.eqn(0)) {
        console.debug(`skipping delegation to ${operator}, amount is zero`)
        continue
      }

      if (stakersBalances.has(owner)) {
        console.log("already has", owner)
        stakersBalances.get(owner).iadd(amount)
      } else {
        stakersBalances.set(owner, amount)
      }
    }

    console.info(
      `found ${stakersBalances.length} holders at block ${this.finalBlock}`
    )

    console.info(
      `dump staked balances to a file: ${TOKEN_STAKING_BALANCES_DUMP_PATH}`
    )
    writeFileSync(
      TOKEN_STAKING_BALANCES_DUMP_PATH,
      JSON.stringify(mapToObject(stakersBalances), null, 2)
    )

    return stakersBalances
  }

  /**
   * @returns {Map<Address,BN>} Token holdings at the final blocks.
   */
  async getTokenHoldingsAtFinalBlock() {
    await this.initialize()

    const allStakeOwners = await this.findHistoricStakeOperatorsOwners()

    const filteredStakeOwnersEOA = await this.filterEOA(allStakeOwners)

    const ownersBalances = await this.checkStakedValues(filteredStakeOwnersEOA)
    return ownersBalances
  }
}
