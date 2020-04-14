import { LedgerSubprovider } from '@0x/subproviders'
import TransportU2F from '@ledgerhq/hw-transport-u2f'
import AppEth from '@ledgerhq/hw-app-eth'
import { AbstractHardwareWalletConnector } from './abstract'
import { getChainIdFromV, getEthereumTxObj } from './utils'
import web3Utils from 'web3-utils'
import { getBufferFromHex } from '../utils/general.utils'

export class LedgerProvider extends AbstractHardwareWalletConnector {
  constructor(chainId) {
    super(new CustomLedgerSubprovider(1101))
  }
}

const ledgerEthereumClientFactoryAsync = async () => {
  const ledgerConnection = await TransportU2F.create()
  ledgerConnection.setExchangeTimeout(100000)
  const ledgerEthClient = new AppEth(ledgerConnection)

  return ledgerEthClient
}

class CustomLedgerSubprovider extends LedgerSubprovider {
  chainId

  constructor(chainId) {
    super({ networkId: chainId, ledgerEthereumClientFactoryAsync, baseDerivationPath: '44\'/60\'' })
    this.chainId = chainId
  }

  async signTransactionAsync(txData) {
    LedgerSubprovider._validateTxParams(txData)
    if (txData.from === undefined || !web3Utils.isAddress(txData.from)) {
      throw new Error('Invalid address')
    }
    txData.chainId = this.chainId

    const initialDerivedKeyInfo = await this._initialDerivedKeyInfoAsync()
    const derivedKeyInfo = this._findDerivedKeyInfoForAddress(initialDerivedKeyInfo, txData.from)
    const fullDerivationPath = derivedKeyInfo.derivationPath

    try {
      this._ledgerClientIfExists = await this._createLedgerClientAsync()
      const tx = getEthereumTxObj(txData, this.chainId)

      tx.raw[6] = this.chainId
      tx.raw[7] = Buffer.from([])
      tx.raw[8] = Buffer.from([])
      const result = await this._ledgerClientIfExists.signTransaction(
        fullDerivationPath,
        tx.serialize().toString('hex')
      )

      let v = result.v
      const rv = parseInt(v, 16)
      let cv = this.chainId * 2 + 35
      if (rv !== cv && (rv & cv) !== rv) {
        cv += 1 // add signature v bit.
      }
      v = cv.toString(16)

      tx.v = getBufferFromHex(v.toString(16))
      tx.r = getBufferFromHex(result.r)
      tx.s = getBufferFromHex(result.s)

      const chainIdFromV = getChainIdFromV(v)
      if (chainIdFromV !== this.chainId) {
        throw new Error('Invalid chainID')
      }

      return `0x${tx.serialize().toString('hex')}`
    } catch (error) {
      await this._destroyLedgerClientAsync()
      throw error
    }
  }
}
