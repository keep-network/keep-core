import { ContractsLoaded } from "../contracts"
const LPRewardsToUniswapTokenAddressCache = {}

// TODO implement functions
export const fetchWrappedTokenBalance = async (address, contractName) => {
  const contracts = await ContractsLoaded
  const lpRewardsContract = contracts[contractName]
  if (!LPRewardsToUniswapTokenAddress[contractName]) {
    LPRewardsToUniswapTokenAddress[
      contractName
    ] = await lpRewardsContract.methods.wrappedToken().call()
  }
}
export const fetchLPRewardsBalance = async (address, contractName) => {}
export const fetchLPRewardsTotalSupply = async (address, contractName) => {}
export const fetchRewardBalance = async (address, contractName) => {}
