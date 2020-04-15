import Common from 'ethereumjs-common'
import { Transaction as EthereumTx } from 'ethereumjs-tx'

export const getEthereumTxObj = (txData, chainId) => {
  const customCommon = Common.forCustomChain('mainnet', {
    name: 'keep-dev',
    chainId,
  })
  const common = new Common(customCommon._chainParams, 'petersburg', ['petersburg'])
  return new EthereumTx(txData, { common })
}

// EIP-155 https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
// v = CHAIN_ID * 2 + 35 => CHAIN_ID = (v - 35) / 2
export const getChainIdFromV = (vInHex) => {
  const vIntValue = parseInt(vInHex, 16)
  const chainId = Math.floor((vIntValue - 35) / 2)
  return chainId < 0 ? 0 : chainId
}

export const getChainId = () => {
  return process.env.CHAIN_ID || 1337
}

export const getRpcURL = () => {
  return process.env.ETH_RPC_URL || 'ws://localhost:8545'
}
