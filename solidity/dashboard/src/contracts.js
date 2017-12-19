import contract from 'truffle-contract'
import Network from "./network"

export async function getKeepToken(address) {
  const KeepToken = contract(require('../../build/contracts/KeepToken.json'))
  const provider = await Network.provider()
  KeepToken.setProvider(provider)
  return KeepToken.at(address)
}

export async function getTokenStaking(address) {
  const TokenStaking = contract(require('../../build/contracts/TokenStaking.json'))
  const provider = await Network.provider()
  TokenStaking.setProvider(provider)
  return TokenStaking.at(address)
}

export async function getTokenVesting(address) {
  const TokenVesting = contract(require('../../build/contracts/TokenVesting.json'))
  const provider = await Network.provider()
  TokenVesting.setProvider(provider)
  return TokenVesting.at(address)
}
