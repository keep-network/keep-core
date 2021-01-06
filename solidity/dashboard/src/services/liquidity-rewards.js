import web3Utils from "web3-utils"
import { Web3Loaded, createERC20Contract } from "../contracts"
import BigNumber from "bignumber.js"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"
import { toTokenUnit } from "../utils/token.utils"
import { getPairData, getKeepTokenPriceInUSD } from "./uniswap-api"

// lp contract address -> wrapped ERC20 token as web3 contract instance
const LPRewardsToWrappedTokenCache = {}
const WEEKS_IN_YEAR = 52

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

export const calculateAPY = async (totalSupplyOfLPRewards, pairSymbol) => {
  totalSupplyOfLPRewards = toTokenUnit(totalSupplyOfLPRewards)

  const pairData = await getPairData(LIQUIDITY_REWARD_PAIRS[pairSymbol].address)

  const rewardPoolPerWeek = new BigNumber(
    LIQUIDITY_REWARD_PAIRS[pairSymbol].rewardPoolPerWeek
  )

  const totalLPTokensInLPRewardsInUSD = totalSupplyOfLPRewards
    .multipliedBy(pairData.reserveUSD)
    .div(pairData.totalSupply)

  const ethPrice = new BigNumber(pairData.reserveUSD).div(pairData.reserveETH)

  let keepTokenInUSD = 0
  if (pairData.token0.symbol === "KEEP") {
    keepTokenInUSD = ethPrice.multipliedBy(pairData.token0.derivedETH)
  } else if (pairData.token1.symbol === "KEEP") {
    keepTokenInUSD = ethPrice.multipliedBy(pairData.token1.derivedETH)
  } else {
    keepTokenInUSD = await getKeepTokenPriceInUSD()
  }

  const r = keepTokenInUSD
    .multipliedBy(rewardPoolPerWeek)
    .div(totalLPTokensInLPRewardsInUSD)

  return r.plus(1).pow(WEEKS_IN_YEAR).minus(1)
}
