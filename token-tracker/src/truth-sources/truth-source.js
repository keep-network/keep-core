/** @typedef { import("../lib/context").default } Context */
/** @typedef { import("../lib/ethereum-helper").Address } Address */
/** @typedef { import("bn.js") } BN */

/**
 * @typedef {Object} ITruthSource
 * @prop {Context} context Description.
 * @prop {Number} targetBlock Description.
 */

import Web3 from "web3"
const { toChecksumAddress } = Web3.utils

/**
 * This is an interface that should be implemented by all sources of truth for
 * token ownership inspection.
 * @property {Context} context Context instance.
 * @property {Number} targetBlock Block at which ownership should be inspected.
 * @property {IgnoredContractsResolver} ignoredContractsResolver Resolver for
 * filtering owners that are contracts. If the owner is one of the ignored contracts
 * it should be removed as it should be already handled by a truth source implementation,
 * to get an actual owner.
 */
export class ITruthSource {
  constructor(
    /** @type {Context} */ context,
    /** @type {Number} */ targetBlock
  ) {
    this.context = context
    this.targetBlock = targetBlock

    this.addressesResolver = new AddressesResolver(context.web3)
  }

  /**
   * Returns a map of addresses with their token holdings discovered using the
   * specific source of truth.
   * @return {Map<Address,BN>} Token holdings at the target blocks.
   */
  async getTokenHoldingsAtTargetBlock() {
    throw new Error("getTokenHoldingsAtTargetBlock not implemented")
  }
}

/**
 * Addresses that are owners of tokens can be externally owned accounts or contracts.
 * Some tokens may be locked in contracts like TokenStaking or TokenGrant, but
 * their actual owner may be another address. Truth of sources should find actual
 * owners of the tokens. This class is used to filter out any addresses that are
 * contracts handled by truth of sources implementations.
 */
export class AddressesResolver {
  /**
   * List of known contracts' addresses that should be ignored as token holders
   * as source of truth implementation resolved actual owners of the tokens locked
   * in these contracts.
   */
  static IGNORED_CONTRACTS = {
    KeepToken: "0x85Eee30c52B0b379b046Fb0F85F4f3Dc3009aFEC",
    OldTokenStaking: "0x6D1140a8c8e6Fac242652F0a5A8171b898c67600",
    TokenStaking: "0x1293a54e160D1cd7075487898d65266081A15458",
    TokenStakingEscrow: "0xDa534b567099Ca481384133bC121D5843F681365",
    StakingPortBacker: "0x236aa50979D5f3De3Bd1Eeb40E81137F22ab794b",
    TokenGrant: "0x175989c71Fd023D580C65F5dC214002687ff88B7",
  }

  /**
   * Each stake from a grant will deploy a separate TokenGrantStake contract instance
   * we will compare contract's bytecode with the reference contract to determine
   * if a contract is an instance of TokenGrantStake.
   */
  static TOKEN_GRANT_STAKE_REF_ADDRESS =
    "0xd966ed532B0F11272D3f5180DB822fB04238D1dE"

  static EOA = "externally owned account"
  static UNKNOWN_CONTRACT = "unknown contract"

  constructor(/** @type {Web3} */ web3) {
    this.web3 = web3
  }

  /**
   * Checks if an address should be ignored when resolving token ownership.
   * It checks the code deployed under the address to determine type of the address.
   * It returns a boolean flag stating if the address should be ignored additionally
   * with a type resolved for an address.
   *
   * It returns following boolean flags for each of the supported address types:
   * - one of ignored contracts addresses: true
   * - externally owned account: false
   * - an instance of TokenGrantStake: true
   * - any other contract: false
   *
   * @param {String} address Address to check.
   * @return {Promise<[Boolean,String]>} Boolean flag if the address should be ignored
   * and resolved type of an address, that can be used for logging on other usage
   * specific for the given type.
   */
  async isIgnoredAddress(address) {
    for (const [ignoredContractName, ignoredContractAddress] of Object.entries(
      AddressesResolver.IGNORED_CONTRACTS
    )) {
      if (
        toChecksumAddress(ignoredContractAddress) === toChecksumAddress(address)
      ) {
        return [true, ignoredContractName]
      }
    }

    const code = await this.web3.eth.getCode(address)

    if (code === "0x") return [false, AddressesResolver.EOA] // address is not a contract

    if (await this.isTokenGrantStake(code)) {
      return [true, "TokenGrantStake"]
    }

    return [false, AddressesResolver.UNKNOWN_CONTRACT]
  }

  /**
   * Checks if a contract under the given address is an instance of TokenGrantStake
   * contract.
   * @param {String} code
   * @return {Boolean} True if it's an instance of TokenGrantStake, false otherwise.
   */
  async isTokenGrantStake(code) {
    if (!this.tokenGrantStakeRefByteCode) {
      this.tokenGrantStakeRefByteCode = await this.web3.eth.getCode(
        AddressesResolver.TOKEN_GRANT_STAKE_REF_ADDRESS
      )
    }

    return code === this.tokenGrantStakeRefByteCode
  }
}
