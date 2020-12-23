import web3Utils from "web3-utils"
import { Web3Loaded, createERC20Contract } from "../contracts"

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

export const fetchLPRewardsTotalSupply = async (LPrewardsContract) => {
  return await LPrewardsContract.methods.totalSupply().call()
}
export const fetchRewardBalance = async (address, LPrewardsContract) => {
  return await LPrewardsContract.methods.earned(address).call()
}
