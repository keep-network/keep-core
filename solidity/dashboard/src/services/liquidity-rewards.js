import web3Utils from "web3-utils"
import { Web3Loaded, createERC20Contract } from "../contracts"
import BigNumber from "bignumber.js"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"

/**
 * Dec 23 2020 15:40 UTC
 *
 * from:
 * https://info.uniswap.org/pair/0xe6f19dab7d43317344282f803f8e8d240708174a
 * https://info.uniswap.org/pair/0x38c8ffee49f286f25d25bad919ff7552e5daf081
 * https://info.uniswap.org/pair/0x854056fd40c1b52037166285b2e54fee774d33f6
 */
const totalLiquidityAmount = {
  KEEP_ETH: new BigNumber(1180514),
  KEEP_TBTC: new BigNumber(113.24),
  TBTC_ETH: new BigNumber(9865),
}

// from https://info.uniswap.org/token/0x85eee30c52b0b379b046fb0f85f4f3dc3009afec (Dec 23 2020 15:34 UTC)
const keepInUSDInBN = new BigNumber(0.23)

// lp contract address -> wrapped ERC20 token as web3 contract instance
const LPRewardsToWrappedTokenCache = {}

export const fetchWrappedTokenBalance = async (address, LPrewardsContract) => {
  const ERC20Contract = await getWrappedTokenConctract(LPrewardsContract)

  return await ERC20Contract.methods.balanceOf(address).call()
}

export const getWrappedTokenConctract = async (LPRewardsContract) => {
  const web3 = await Web3Loaded
  const lpRewardsContractAddress = web3Utils.toChecksumAddress(
    LPRewardsContract.options.address
  )

  if (!LPRewardsToWrappedTokenCache.hasOwnProperty(lpRewardsContractAddress)) {
    const wrappedTokenAddress = await LPRewardsContract.methods
      .wrappedToken()
      .call()
    LPRewardsToWrappedTokenCache[
      lpRewardsContractAddress
    ] = createERC20Contract(web3, wrappedTokenAddress)
  }

  return LPRewardsToWrappedTokenCache[lpRewardsContractAddress]
}

export const fetchStakedBalance = async (address, LPrewardsContract) => {
  return await LPrewardsContract.methods.balanceOf(address).call()
}

export const fetchTotalLPTokensCreatedInUniswap = async (LPrewardsContract) => {
  const ERC20Contract = await getWrappedTokenConctract(LPrewardsContract)
  return await ERC20Contract.methods.totalSupply().call()
}

export const fetchLPRewardsTotalSupply = async (LPrewardsContract) => {
  return await LPrewardsContract.methods.totalSupply().call()
}
export const fetchRewardBalance = async (address, LPrewardsContract) => {
  return await LPrewardsContract.methods.earned(address).call()
}

export const calculateAPY = (
  totalSupplyInWei,
  totalLPTokensCreatedInUniswapInWei,
  pairSymbol
) => {
  const totalSupply = web3Utils.fromWei(totalSupplyInWei)
  const totalLPTokensCreatedInUniswap = web3Utils.fromWei(
    totalLPTokensCreatedInUniswapInWei
  )
  const totalSupplyInBN = new BigNumber(totalSupply)
  const totalLPTokensCreatedInUniswapInBN = new BigNumber(
    totalLPTokensCreatedInUniswap
  )
  // TODO : Add pairSymbol to constants (probably not be needed)
  const rewardPoolPerWeekInBN = new BigNumber(
    LIQUIDITY_REWARD_PAIRS[pairSymbol].rewardPoolPerWeek
  )

  const totalLPTokensInLPRewardsInUSD = totalSupplyInBN
    .multipliedBy(totalLiquidityAmount[pairSymbol])
    .div(totalLPTokensCreatedInUniswapInBN)

  const r = keepInUSDInBN
    .multipliedBy(rewardPoolPerWeekInBN)
    .div(totalLPTokensInLPRewardsInUSD)

  // TODO: Add 52 to constants
  const apy = r.plus(1).pow(52).minus(1)

  return apy
}
