import { LedgerSubprovider } from '@0x/subproviders'
import TransportU2F from '@ledgerhq/hw-transport-u2f'
import AppEth from '@ledgerhq/hw-app-eth'
import { AbstractHardwareWalletConnector } from './abstract-connector'
import { getChainIdFromV, getEthereumTxObj, getChainId } from './utils'
import web3Utils from 'web3-utils'
import { getBufferFromHex } from '../utils/general.utils'

export class LedgerProvider extends AbstractHardwareWalletConnector {
  constructor() {
    super(new CustomLedgerSubprovider(getChainId()))
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

      tx.raw[6] = getBufferFromHex(this.chainId.toString(16))
      tx.raw[7] = Buffer.from([])
      tx.raw[8] = Buffer.from([])
      const result = await this._ledgerClientIfExists.signTransaction(
        fullDerivationPath,
        tx.serialize().toString('hex')
      )

      // The transport layer only returns the lower 2 bytes.
      // The returned `v` will be wrong for chainId's < 255 and has to be recomputed.
      const ledgerSignedV = parseInt(result.v, 16)
      let signedV = this.chainId * 2 + 35
      if (ledgerSignedV % 2 === 0) {
        signedV += 1
      }

      tx.v = getBufferFromHex(signedV.toString(16))
      tx.r = getBufferFromHex(result.r)
      tx.s = getBufferFromHex(result.s)

      // Compare `v` value returned from Ledger.
      // eg. for `chainId = 1101` => `v = 2238(08be)` => `ledgerSignedV = 190(be)`
      // `2238 & 0xff = 190`
      const isValidSignedV = (signedV & 0xff) === ledgerSignedV
      const chainIdFromV = getChainIdFromV(tx.v)
      if ((chainIdFromV !== this.chainId) && !isValidSignedV) {
        throw new Error('Invalid chainID')
      }

      return `0x${tx.serialize().toString('hex')}`
    } catch (error) {
      await this._destroyLedgerClientAsync()
      throw error
    }
  }
}
