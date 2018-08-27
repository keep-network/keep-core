import contract from 'truffle-contract'
import Network from "./network"

export async function getKeepToken(address) {

  const code = await Network.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const KeepToken = contract(require('contracts/KeepToken.json'))
  const provider = await Network.provider()
  KeepToken.setProvider(provider)
  return KeepToken.at(address)
}

export async function getTokenStaking(address) {

  const code = await Network.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const TokenStaking = contract(require('contracts/TokenStaking.json'))
  const provider = await Network.provider()
  TokenStaking.setProvider(provider)
  return TokenStaking.at(address)
}

export async function getTokenGrant(address) {

  const code = await Network.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const TokenGrant = contract(require('contracts/TokenGrant.json'))
  const provider = await Network.provider()
  TokenGrant.setProvider(provider)
  return TokenGrant.at(address)
}
