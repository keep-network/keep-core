import TrezorConnect from 'trezor-connect'
import { TrezorSubprovider } from '@0x/subproviders/lib/src/subproviders/trezor'
import { AbstractHardwareWalletConnector } from './abstract'
import web3Utils from 'web3-utils'
import EthereumTx from 'ethereumjs-tx'
import Common from 'ethereumjs-common'

TrezorConnect.init({
  lazyLoad: true, // this param will prevent iframe injection until TrezorConnect.method will be called
  manifest: {
    email: 'work@keep.network',
    appUrl: 'keep.network',
  },
  popup: true,
  //  debug: true,
})

export class TrezorProvider extends AbstractHardwareWalletConnector {
  constructor() {
    super(new Trezor({
      trezorConnectClientApi: TrezorConnect,
      networkId: 1337,
    }))
  }
}

class Trezor extends TrezorSubprovider {
  constructor(config) {
    super(config)
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
        chainId: this._networkId,
        nonce: txData.nonce,
        gasLimit: txData.gas,
        gasPrice: txData.gasPrice,
      },
    })
    if (response.success) {
      const payload = response.payload
      const customCommon = Common.forCustomChain('mainnet', {
        name: 'keep-dev',
        chainId: this._networkId,
      })
      const tx = new EthereumTx(txData, { common: customCommon })
      tx.v = getBufferFromHex(payload.v)
      tx.r = getBufferFromHex(payload.r)
      tx.s = getBufferFromHex(payload.s)

      return `0x${tx.serialize().toString('hex')}`
    } else {
      const payload = response.payload
      throw new Error(payload.error)
    }
  }
}

const getBufferFromHex = (hex) => {
  hex = hex.substring(0, 2) == '0x' ? hex.substring(2) : hex
  if (hex == '') {
    return new Buffer('', 'hex')
  }
  const padLeft = hex.length % 2 != 0 ? '0' + hex : hex
  return new Buffer(padLeft.toLowerCase(), 'hex')
}
