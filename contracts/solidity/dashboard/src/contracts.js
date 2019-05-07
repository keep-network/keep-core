import contract from 'truffle-contract'

export async function getKeepToken(web3, address) {

  const code = await web3.eth.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const KeepToken = contract(require('contracts/KeepToken.json'))
  KeepToken.setProvider(web3.currentProvider)
  return KeepToken.at(address)
}

export async function getTokenStaking(web3, address) {

  const code = await web3.eth.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const Staking = contract(require('contracts/Staking.json'))
  Staking.setProvider(web3.currentProvider)
  return Staking.at(address)
}

export async function getTokenGrant(web3, address) {

  const code = await web3.eth.getCode(address);
  if (!code || code === '0x0' || code === '0x') throw Error('No contract at address');

  const TokenGrant = contract(require('contracts/TokenGrant.json'))
  TokenGrant.setProvider(web3.currentProvider)
  return TokenGrant.at(address)
}
