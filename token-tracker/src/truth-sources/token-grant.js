/** @typedef { import("../lib/context.js").default } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */
/** @typedef { import("../lib/contract-helper").ContractInstance } ContractInstance*/
/** @typedef { import("bn.js") } BN */

import { ITruthSource } from "./truth-source.js"
import { Contract } from "../lib/contract-helper.js"
import { logger } from "../lib/winston.js"

import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"

import TokenGrantJSON from "@keep-network/keep-core/artifacts/TokenGrant.json"
import TokenStakingEscrowJSON from "@keep-network/keep-core/artifacts/TokenStakingEscrow.json"

import Web3 from "web3"
const { toBN } = Web3.utils

import { resolveGrantee } from "../lib/owner-lookup.js"

const TOKEN_GRANT_OUTPUT_PATH = "./tmp/token-grant.json"
const TOKEN_GRANT_BALANCES_OUTPUT_PATH = "./tmp/token-grant-balances.json"

export class TokenGrantTruthSource extends ITruthSource {
  /** @property {ContractInstance} tokenGrant */
  /** @property {ContractInstance} tokenStaking */
  /** @property {ContractInstance} oldTokenStaking */
  /** @property {ContractInstance} tokenStakingEscrow */

  constructor(
    /** @type {Context} */ context,
    /** @type {Number}*/ targetBlock
  ) {
    super(context, targetBlock)
  }

  async initialize() {
    this.context.addContract(
      "TokenGrant",
      new Contract(TokenGrantJSON, this.context.web3)
    )
    this.context.addContract(
      "TokenStakingEscrow",
      new Contract(TokenStakingEscrowJSON, this.context.web3)
    )

    this.tokenGrant = await this.context.getContract("TokenGrant").deployed()
    this.tokenStaking = await this.context
      .getContract("TokenStaking")
      .deployed()
    this.oldTokenStaking = await this.context
      .getContract("OldTokenStaking")
      .deployed()
    this.tokenStakingEscrow = await this.context
      .getContract("TokenStakingEscrow")
      .deployed()
  }

  /**
   * @typedef {Object} TokenGrant
   * @property {string} id
   * @property {Address} grantee
   * @property {BN} amount
   * @property {BN} withdrawn
   * @property {BN} revoked
   * @property {BN} slashed
   * @property {BN} seized
   * @property {BN} escrowWithdrawn
   * @property {BN} escrowRevoked
   * @property {BN} balance
   */

  /**
   * Gets all token grants based on TokenGrantCreated event of TokenGrant contract.
   * Returns details of each grant like ID, grantee and balance. The balance is
   * calculated at the pre defined target block as total grant amount minus withdrawn
   * and revoked tokens. The balance includes tokens that were staked, it's important
   * to exclude the stake amount in TokenStaking source of truth to avoid value
   * duplication.
   * @return {TokenGrant[]} Token grant with ID, grantee and balance at target block.
   */
  async findTokenGrants() {
    logger.info(
      `looking for TokenGrantCreated events emitted from ${this.tokenGrant.options.address} ` +
        `between blocks ${this.context.deploymentBlock} and ${this.targetBlock}`
    )

    /** @type {Array} */
    const events = await getPastEvents(
      this.context.web3,
      this.tokenGrant,
      "TokenGrantCreated",
      this.context.deploymentBlock,
      this.targetBlock
    )
    logger.info(`found ${events.length} token grant created events`)

    /** @type {String[]} */
    const tokenGrantIDs = events.map((event) => event.returnValues.id)

    /** @type {TokenGrant[]} */
    const tokenGrants = []
    for (const id of tokenGrantIDs) {
      logger.debug(`processing grant with id: ${id}`)

      const grantDetails = await callWithRetry(
        this.tokenGrant.methods.getGrant(id),
        undefined,
        undefined,
        this.targetBlock
      )

      /** @type {TokenGrant} */
      const grant = {
        id: id,
        grantee: grantDetails.grantee,
        amount: toBN(grantDetails.amount),
        withdrawn: toBN(grantDetails.withdrawn),
        revoked: toBN(grantDetails.revokedAmount),
      }

      // Balance is the total grant amount minus revoked and withdrawn tokens. The
      // value includes staked tokens.
      grant.balance = grant.amount.sub(grant.withdrawn).sub(grant.revoked)

      // If tokens were staked.
      if (toBN(grantDetails.staked).gtn(0)) {
        const operators = await this.getOperatorsForGrant(id)

        // Check for any seized or slashed ones.
        ;[grant.slashed, grant.seized] = await this.getSlashedSeizedAmount(
          operators
        )

        if (grant.slashed.gtn(0))
          logger.debug(`grant ${id}: slashed tokens: ${grant.slashed}`)
        if (grant.seized.gtn(0))
          logger.debug(`grant ${id}: seized tokens: ${grant.seized}`)

        grant.balance.isub(grant.slashed)
        grant.balance.isub(grant.seized)

        // Check for any withdrawn or revoked tokens in TokenStakingEscrow.
        ;[
          grant.escrowWithdrawn,
          grant.escrowRevoked,
        ] = await this.getWithdrawnRevokedAmountFromEscrow(operators)

        if (grant.escrowWithdrawn.gtn(0))
          logger.debug(
            `grant ${id}: escrow withdrawn tokens: ${grant.escrowWithdrawn}`
          )
        if (grant.escrowRevoked.gtn(0))
          logger.debug(
            `grant ${id}: escrow revoked tokens: ${grant.escrowRevoked}`
          )

        grant.balance.isub(grant.escrowWithdrawn)
        grant.balance.isub(grant.escrowRevoked)
      }

      tokenGrants.push(grant)
    }

    dumpDataToFile(tokenGrants, TOKEN_GRANT_OUTPUT_PATH)

    return tokenGrants
  }

  /**
   * @param {String} grantID
   * @return {Address[]}
   */
  async getOperatorsForGrant(grantID) {
    /** @type {Set<Address>} */
    const operators = new Set()

    /** @type {Array} */
    const grantStakedEvents = await getPastEvents(
      this.context.web3,
      this.tokenGrant,
      "TokenGrantStaked",
      this.context.deploymentBlock,
      this.targetBlock,
      { grantId: grantID }
    )
    logger.debug(
      `found ${grantStakedEvents.length} TokenGrantStaked events for grant ${grantID}`
    )

    grantStakedEvents.forEach((event) =>
      operators.add(event.returnValues.operator)
    )

    /** @type {Array} */
    const tokenStakingEscrowEvents = await getPastEvents(
      this.context.web3,
      this.tokenStakingEscrow,
      "DepositRedelegated",
      this.context.deploymentBlock,
      this.targetBlock,
      { grantId: grantID }
    )
    logger.debug(
      `found ${tokenStakingEscrowEvents.length} DepositRedelegated events for grant ${grantID}`
    )

    tokenStakingEscrowEvents.forEach((event) =>
      operators.add(event.returnValues.newOperator)
    )

    return Array.from(operators)
  }

  /**
   * @param {Address[]} operators
   * @return {Promise<[BN,BN]>}
   */
  async getWithdrawnRevokedAmountFromEscrow(operators) {
    const withdrawnAmount = toBN(0)
    const revokedAmount = toBN(0)

    for (const operator of operators) {
      logger.debug(
        `looking for tokens withdrawn or revoked from escrow for operator: ${operator}`
      )

      /** @type {Array} */
      const withdrawnEvents = await getPastEvents(
        this.context.web3,
        this.tokenStakingEscrow,
        "DepositWithdrawn",
        this.context.deploymentBlock,
        this.targetBlock,
        { operator: operator }
      )

      withdrawnEvents.forEach((event) =>
        withdrawnAmount.iadd(toBN(event.returnValues.amount))
      )

      /** @type {Array} */
      const revokedEvents = await getPastEvents(
        this.context.web3,
        this.tokenStakingEscrow,
        "RevokedDepositWithdrawn",
        this.context.deploymentBlock,
        this.targetBlock,
        { operator: operator }
      )

      revokedEvents.forEach((event) =>
        revokedAmount.iadd(toBN(event.returnValues.amount))
      )
    }

    return [withdrawnAmount, revokedAmount]
  }

  /**
   * @param {Address[]} operators
   * @return {Promise<[BN,BN]>}
   */
  async getSlashedSeizedAmount(operators) {
    const slashedAmount = toBN(0)
    const seizedAmount = toBN(0)
    for (const operator of operators) {
      logger.debug(
        `looking for slashed or seized tokens for operator: ${operator}`
      )

      slashedAmount.iadd(
        await this.getAmountFromTokenStakingEvent("TokensSlashed", operator)
      )
      seizedAmount.iadd(
        await this.getAmountFromTokenStakingEvent("TokensSeized", operator)
      )
    }

    return [slashedAmount, seizedAmount]
  }

  /**
   * @param {String} eventName Name of the event to check.
   * @param {Address} operator Operator address to filter events.
   * @return {Promise<BN>}
   */
  async getAmountFromTokenStakingEvent(eventName, operator) {
    /** @type {Array} */
    const eventsTokenStaking = await getPastEvents(
      this.context.web3,
      this.tokenStaking,
      eventName,
      this.context.deploymentBlock,
      this.targetBlock,
      { operator: operator }
    )

    /** @type {Array} */
    const eventsOldTokenStaking = await getPastEvents(
      this.context.web3,
      this.oldTokenStaking,
      eventName,
      this.context.deploymentBlock,
      this.targetBlock,
      { operator: operator }
    )

    const totalAmount = toBN(0)

    eventsTokenStaking.forEach((event) =>
      totalAmount.iadd(toBN(event.returnValues.amount))
    )
    eventsOldTokenStaking.forEach((event) =>
      totalAmount.iadd(toBN(event.returnValues.amount))
    )

    return toBN(totalAmount)
  }

  /**
   * Converts list of grants to a map of owners and their balances. It resolves
   * owner for each grantee and combines balances for the same owners.
   * @param {Array<TokenGrant>} tokenGrants Stake owners to filter.
   * @return {Map<Address,BN>} Map of owners and their balances.
   */
  async resolveOwnersBalances(tokenGrants) {
    /** @type {Map<Address,BN>} */
    const ownersBalances = new Map()

    for (const tokenGrant of tokenGrants) {
      const { grantee, balance } = tokenGrant

      const owner = await resolveGrantee(this.context.web3, grantee)

      const { toChecksumAddress } = this.context.web3.utils

      if (toChecksumAddress(owner) !== toChecksumAddress(grantee))
        logger.debug(
          `resolved grantee for managed grant [${grantee}]: ${owner}`
        )

      if (ownersBalances.has(owner)) {
        ownersBalances.get(owner).iadd(balance)
      } else {
        ownersBalances.set(owner, balance)
      }
    }

    dumpDataToFile(ownersBalances, TOKEN_GRANT_BALANCES_OUTPUT_PATH)

    return ownersBalances
  }

  /**
   * Returns a map of addresses with their token holdings based on Token Grants.
   * @return {Map<Address,BN>} Token holdings at the target block.
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const allGrants = await this.findTokenGrants()

    const ownersBalances = await this.resolveOwnersBalances(allGrants)

    return ownersBalances
  }
}
