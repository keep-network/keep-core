import axios from "axios"
import BigNumber from "bignumber.js"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"

const UNISWAP_API_URL =
  "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"

export const getPairData = async (pairId) => {
  const response = await axios.post(UNISWAP_API_URL, {
    query: `query pairquery {
        pair(id: "${pairId}") {
            token0 {
                symbol,
                derivedETH,
            },
            token1 {
                symbol,
                derivedETH,
            },
            reserve0,
            reserve1,
            reserveETH,
            reserveUSD,
            token0Price,
            token1Price,
            totalSupply,
        }
    }`,
  })

  if (response.data?.errors) {
    const error = new Error(
      "Failed fetching data from Uniswap V2 subgraph API."
    )
    error.errors = response.data.errors
    throw error
  }

  return response.data.data.pair
}

const getTokenPriceInUSD = async (address) => {
  const pairData = await getPairData(address)
  const ethPrice = new BigNumber(pairData.reserveUSD).div(pairData.reserveETH)

  return ethPrice.multipliedBy(pairData.token0.derivedETH)
}

/**
 * Returns the current KEEP token price in USD based on the Uniswap pool.
 *
 * @return {BigNumber} KEEP token price in USD.
 */
export const getKeepTokenPriceInUSD = async () => {
  return await getTokenPriceInUSD(LIQUIDITY_REWARD_PAIRS.KEEP_ETH.address)
}

/**
 * Returns the current BTC price in USD based on the TBTC/ETH Uniswap pool.
 *
 * @return {BigNumber} BTC price in USD.
 */
export const getBTCPriceInUSD = async () => {
  return await getTokenPriceInUSD(LIQUIDITY_REWARD_PAIRS.TBTC_ETH.address)
}
