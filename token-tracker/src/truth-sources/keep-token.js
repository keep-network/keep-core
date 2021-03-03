/** @typedef { import("../lib/context").Context } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */

import BN from "bn.js"
import pAll from "p-all"

import { ITruthSource, AddressesResolver } from "./truth-source.js"

import { Contract } from "../lib/contract-helper.js"

import { EthereumHelpers } from "@keep-network/tbtc.js"
const { callWithRetry } = EthereumHelpers

import { getPastEvents } from "../lib/ethereum-helper.js"
import { dumpDataToFile } from "../lib/file-helper.js"
import { logger } from "../lib/winston.js"

import KeepTokenJSON from "@keep-network/keep-core/artifacts/KeepToken.json"

const KEEP_TOKEN_HISTORIC_HOLDERS_DUMP_PATH = "./tmp/keep-token-holders.json"
const KEEP_TOKEN_BALANCES_DUMP_PATH = "./tmp/keep-token-balances.json"
const KEEP_TOKEN_UNKNOWN_HOLDERS_CONTRACTS_DUMP_PATH =
  "./tmp/keep-token-unknown_holders_contracts.json"

const CONCURRENCY_LEVEL = 20

export class KeepTokenTruthSource extends ITruthSource {
  /**
   * @param {Context} context
   * @param {Number} targetBlock
   */
  constructor(context, targetBlock) {
    super(context, targetBlock)
  }

  async initialize() {
    this.context.addContract(
      "KeepToken",
      new Contract(KeepTokenJSON, this.context.web3)
    )

    this.keepToken = await this.context.getContract("KeepToken").deployed()
  }

  /**
   * Finds all historic holders of KEEP token based on Transfer events.
   * @return {Set<Address>} All historic token holders.
   * */
  async findHistoricHolders() {
    logger.info(
      `looking for Transfer events emitted from ${this.keepToken.options.address} ` +
        `between blocks ${this.context.deploymentBlock} and ${this.targetBlock}`
    )

    const events = await getPastEvents(
      this.context.web3,
      this.keepToken,
      "Transfer",
      this.context.deploymentBlock,
      this.targetBlock
    )
    logger.info(`found ${events.length} token transfer events`)

    const tokenHolders = new Set()
    events.forEach((event) => tokenHolders.add(event.returnValues.to))

    logger.info(`found ${tokenHolders.size} unique historic holders`)

    dumpDataToFile(tokenHolders, KEEP_TOKEN_HISTORIC_HOLDERS_DUMP_PATH)

    return tokenHolders
  }

  /**
   * Filters addresses based on the rules defined for sources of truth. Ignores
   * addresses of known Keep contracts for which actual holders are resolved.
   * @param {Set<Address>} addresses Token holders addresses.
   * @return {Set<Address>} Filtered set of addresses.
   */
  async filterHolders(addresses) {
    logger.info(`filter holders`)

    const unknownContracts = new Set()

    /**
     *
     * @param {AddressesResolver} addressesResolver
     * @param {Address} address
     */
    const checkAddress = async (addressesResolver, address) => {
      const [isIgnored, addressType] = await addressesResolver.isIgnoredAddress(
        address
      )

      logger.debug(`${address} is ${addressType}`)

      if (isIgnored) {
        return false
      }

      // If the address doesn't match any of ignored contracts store it for
      // reference if we want to double check them.
      if (addressType === AddressesResolver.UNKNOWN_CONTRACT)
        unknownContracts.add(address)

      return true
    }

    // TODO: Consider moving addresses ignore out from soure of truth scripts to
    // the very final step after we've got a combined map of holdings.
    const concurrentAddressChecks = Array.from(addresses).map((address) => () =>
      new Promise(async (resolve, reject) => {
        await checkAddress(this.addressesResolver, address)
          .then((result) => {
            if (result) {
              resolve(address)
            } else {
              resolve()
            }
          })
          .catch(reject)
      })
    )

    const filteredHolders = await pAll(concurrentAddressChecks, {
      concurrency: CONCURRENCY_LEVEL,
    })

    dumpDataToFile(
      unknownContracts,
      KEEP_TOKEN_UNKNOWN_HOLDERS_CONTRACTS_DUMP_PATH
    )

    return new Set(filteredHolders.filter(Boolean)) // filter out empty entries
  }

  /**
   * Checks token balance for addresses on the list.
   * @param {Set<Address>} tokenHolders Token holders to check.
   * @return {Map<Address,BN>} Token holdings at the target blocks.
   */
  async checkTargetHoldersBalances(tokenHolders) {
    logger.info(`check token holdings at block ${this.targetBlock}`)

    const concurrentBalanceChecks = Array.from(tokenHolders).map(
      (holder) => () =>
        new Promise(async (resolve, reject) => {
          callWithRetry(
            this.keepToken.methods.balanceOf(holder),
            undefined,
            undefined,
            this.targetBlock
          )
            .then((targetBalance) => {
              logger.debug(`holder ${holder} balance ${targetBalance}`)
              resolve({ holder: holder, balance: targetBalance })
            })
            .catch(reject)
        })
    )

    let balances
    try {
      balances = await pAll(concurrentBalanceChecks, {
        concurrency: CONCURRENCY_LEVEL,
      })
    } catch (err) {
      throw new Error(`concurrent execution failed: ${err}`)
    }

    if (tokenHolders.size != balances.length) {
      throw new Error(
        `unexpected number of fetched holders balances; ` +
          `expected: ${tokenHolders.size}, actual: ${balances.length}`
      )
    }

    /** @type {Map<Address,BN>} */
    const holdersBalances = new Map()

    for (const entry of balances) {
      const holder = entry.holder
      const targetBalance = new BN(entry.balance)

      if (targetBalance.gtn(0)) {
        holdersBalances.set(holder, targetBalance)
      }
    }

    logger.info(
      `found ${holdersBalances.size} holders at block ${this.targetBlock}`
    )

    dumpDataToFile(holdersBalances, KEEP_TOKEN_BALANCES_DUMP_PATH)

    return holdersBalances
  }

  /**
   * @return {Map<Address,BN>} Token holdings at the target blocks.
   */
  async getTokenHoldingsAtTargetBlock() {
    await this.initialize()

    const allTokenHolders = await this.findHistoricHolders()

    const filteredTokenHolders = await this.filterHolders(allTokenHolders)

    const holdersBalances = await this.checkTargetHoldersBalances(
      filteredTokenHolders
    )

    return holdersBalances
  }
}
