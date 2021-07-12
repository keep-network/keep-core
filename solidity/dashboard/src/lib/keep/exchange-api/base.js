/** @typedef { import("bignumber.js").BigNumber} BigNumber */

class BaseExchange {
  /**
   * Returns the Uniswap pair data.
   *
   * @param {string} pairId Uniswap pair id- address of the `UniswapV2Pair`
   * contract.
   * @return {Object} Uniswap pair data.
   */
  getUniswapPairData = async (pairId) => {
    return await this._getUniswapPairData(pairId)
  }

  /**
   * Returns the current KEEP token price in USD based on the exchange data.
   *
   * @return {BigNumber} KEEP token price in USD.
   */
  getKeepTokenPriceInUSD = async () => {
    return await this._getKeepTokenPriceInUSD()
  }

  /**
   * Returns the current BTC price in USD based on the exchange data.
   *
   * @return {BigNumber} BTC price in USD.
   */
  getBTCPriceInUSD = async () => {
    return await this._getBTCPriceInUSD()
  }
}

export default BaseExchange
