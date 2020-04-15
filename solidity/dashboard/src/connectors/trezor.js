import TrezorConnect from 'trezor-connect'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import web3Utils from 'web3-utils'
import { AbstractHardwareWalletConnector } from './abstract'
import { getEthereumTxObj, getChainIdFromV } from './utils'
import { getBufferFromHex } from '../utils/general.utils'
import { getChainId } from './utils'

export class TrezorProvider extends AbstractHardwareWalletConnector {
  constructor() {
    super(new CustomTrezorSubprovider(getChainId()))
  }
}

class CustomTrezorSubprovider extends TrezorSubprovider {
  chainId

  constructor(chainId) {
    super({ trezorConnectClientApi: TrezorConnect, networkId: chainId })
    this.chainId = chainId
    this._trezorConnectClientApi.manifest({
      email: 'work@keep.network',
      appUrl: 'https://keep.network',
    })
  }

  async signTransactionAsync(txData) {
    if (txData.from === undefined || !web3Utils.isAddress(txData.from)) {
      throw new Error('Invalid address')
    }
    txData.value = txData.value ? txData.value : '0x0'
    txData.data = txData.data ? txData.data : '0x'
    txData.gas = txData.gas ? txData.gas : '0x0'
    txData.gasPrice = txData.gasPrice ? txData.gasPrice : '0x0'

    const initialDerivedKeyInfo = await this._initialDerivedKeyInfoAsync()
    const derivedKeyInfo = this._findDerivedKeyInfoForAddress(initialDerivedKeyInfo, txData.from)
    const fullDerivationPath = derivedKeyInfo.derivationPath

    const response = await this._trezorConnectClientApi.ethereumSignTransaction({
      path: fullDerivationPath,
      transaction: {
        to: txData.to,
        value: txData.value,
        data: txData.data,
        chainId: this.chainId,
        nonce: txData.nonce,
        gasLimit: txData.gas,
        gasPrice: txData.gasPrice,
      },
    })
    if (!response.success) {
      throw new Error(response.payload.error)
    }
    const { payload: { v, r, s } } = response
    const tx = getEthereumTxObj(txData, this.chainId)

    tx.v = getBufferFromHex(v)
    tx.r = getBufferFromHex(r)
    tx.s = getBufferFromHex(s)
    const chainIdFromV = getChainIdFromV(v)
    if (chainIdFromV !== this.chainId) {
      throw new Error('Invalid chainID')
    }

    return `0x${tx.serialize().toString('hex')}`
  }
}
