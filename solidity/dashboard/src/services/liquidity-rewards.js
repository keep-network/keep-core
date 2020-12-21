import { Web3Loaded, createERC20Contract } from "../contracts"

// lp contract address -> wrapped ERC20 token as web3 contract instance
const LPRewardsToWrappedTokenCache = {}

// TODO implement functions
export const fetchWrappedTokenBalance = async (address, LPrewardsContract) => {
  const web3 = await Web3Loaded
  const { address: lpRewardsContractAddress } = LPrewardsContract.options

  if (!LPRewardsToWrappedTokenCache[lpRewardsContractAddress]) {
    const wrappedTokenAddress = await LPrewardsContract.methods
      .wrappedToken()
      .call()
    LPRewardsToWrappedTokenCache[
      lpRewardsContractAddress
    ] = createERC20Contract(web3, wrappedTokenAddress)
  }

  const ERC20Contract = LPRewardsToWrappedTokenCache[lpRewardsContractAddress]

  return await ERC20Contract.methods.balanceOf(address).call()
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
