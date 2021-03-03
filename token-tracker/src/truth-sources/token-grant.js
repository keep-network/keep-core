/** @typedef { import("../lib/context.js").default } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */
/** @typedef { import("../lib/contract-helper").ContractInstance } ContractInstance*/

import { ITruthSource } from "./truth-source.js"
import { Contract } from "../lib/contract-helper.js"
import { logger } from "../lib/winston.js"

import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"

import TokenGrantJSON from "@keep-network/keep-core/artifacts/TokenGrant.json"

import Web3 from "web3"
const { toBN } = Web3.utils

import { resolveGrantee } from "../lib/owner-lookup.js"

const TOKEN_GRANT_OUTPUT_PATH = "./tmp/token-grant.json"
const TOKEN_GRANT_BALANCES_OUTPUT_PATH = "./tmp/token-grant-balances.json"

export class TokenGrantTruthSource extends ITruthSource {
  /** @property {ContractInstance} tokenGrant */

  constructor(
    /** @type {Context} */ context,
    /** @type {Number}*/ targetBlock
  ) {
    super(context, targetBlock)
  }

  async initialize() {
    const TokenGrant = new Contract(TokenGrantJSON, this.context.web3)

    this.context.addContract("TokenGrant", TokenGrant)

    this.tokenGrant = await this.context.getContract("TokenGrant").deployed()
  }

  /**
   * @typedef {Object} TokenGrant
   * @property {string} id
   * @property {Address} grantee
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
      const grant = await callWithRetry(
        this.tokenGrant.methods.getGrant(id),
        undefined,
        undefined,
        this.targetBlock
      )

      // Balance is the total grant amount minus revoked and withdrawn tokens. The
      // value includes staked tokens.
      const balance = toBN(grant.amount)
        .sub(toBN(grant.withdrawn))
        .sub(toBN(grant.revokedAmount))

      tokenGrants.push({ id: id, grantee: grant.grantee, balance: balance })
    }

    dumpDataToFile(tokenGrants, TOKEN_GRANT_OUTPUT_PATH)

    return tokenGrants
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

      if (owner !== grantee)
        logger.debug(
          `resolved grantee for managed grant [${grantee}]: ${owner}`
        )

      if (ownersBalances.has(owner)) {
        ownersBalances.get(owner).iadd(balance)
      } else {
        ownersBalances.set(owner, toBN(0))
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
