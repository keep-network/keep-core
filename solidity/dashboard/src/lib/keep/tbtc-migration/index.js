/** @typedef { import("../../web3").BaseContract} BaseContract */
/** @typedef { import("../../web3").Web3LibWrapper} Web3LibWrapper */

import { div, mul } from "../../../utils/arithmetics.utils"
import { Token } from "../../../utils/token.utils"

class TBTCV2Migration {
  /**
   * @param {BaseContract} _tbtcV1
   * @param {BaseContract} _tbtcV2
   * @param {BaseContract} _vendingMachine
   * @param {Web3LibWrapper} _web3
   */
  constructor(_tbtcV1, _tbtcV2, _vendingMachine, _web3) {
    this.tbtcV1 = _tbtcV1
    this.tbtcV2 = _tbtcV2
    this.vendingMachine = _vendingMachine
    this.web3 = _web3
  }

  unmintFee = async () => {
    return await this.vendingMachine.makeCall("unmintFee")
  }

  unmintFeeFor = async (amount, unmintFee = null) => {
    if (!unmintFee) {
      return await this.vendingMachine.makeCall("unmintFeeFor", amount)
    }

    return div(
      mul(amount, unmintFee),
      Token.fromTokenUnit(1, 18).toString()
    ).toString()
  }

  tbtcV1BalanceOf = async (address) => {
    return await this.tbtcV1.makeCall("balanceOf", address)
  }

  tbtcV2BalanceOf = async (address) => {
    return await this.tbtcV2.makeCall("balanceOf", address)
  }
}

export default TBTCV2Migration
