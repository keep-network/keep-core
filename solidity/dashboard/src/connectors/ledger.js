import { LedgerSubprovider } from "@0x/subproviders"
import TransportU2F from "@ledgerhq/hw-transport-u2f"
import AppEth from "@ledgerhq/hw-app-eth"
import { AbstractHardwareWalletConnector } from "./abstract-connector"
import { getChainIdFromV, getEthereumTxObj, getChainId } from "./utils"
import web3Utils from "web3-utils"
import { getBufferFromHex } from "../utils/general.utils"

export const LEDGER_DERIVATION_PATHS = {
  LEDGER_LIVE: `m/44'/60'/x'/0/0`,
  LEDGER_LEGACY: `m/44'/60'/0'/x`,
}

export class LedgerProvider extends AbstractHardwareWalletConnector {
  constructor(baseDerivationPath) {
    super(new CustomLedgerSubprovider(getChainId(), baseDerivationPath))
  }
}

const LEDGER_EXCHANGE_TIMEOUT = 100000
const ledgerEthereumClientFactoryAsync = async () => {
  const ledgerConnection = await TransportU2F.create()
  ledgerConnection.setExchangeTimeout(LEDGER_EXCHANGE_TIMEOUT)
  const ledgerEthClient = new AppEth(ledgerConnection)

  return ledgerEthClient
}

class CustomLedgerSubprovider extends LedgerSubprovider {
  chainId
  addressToPathMap = {}
  pathToAddressMap = {}

  constructor(chainId, baseDerivationPath) {
    super({
      networkId: chainId,
      ledgerEthereumClientFactoryAsync,
      baseDerivationPath: baseDerivationPath,
    })
    this.chainId = chainId
  }

  async signTransactionAsync(txData) {
    LedgerSubprovider._validateTxParams(txData)
    if (txData.from === undefined || !web3Utils.isAddress(txData.from)) {
      throw new Error("Invalid address")
    }
    txData.chainId = this.chainId

    const fullDerivationPath = this.addressToPathMap[
      web3Utils.toChecksumAddress(txData.from)
    ]

    try {
      this._ledgerClientIfExists = await this._createLedgerClientAsync()
      const tx = getEthereumTxObj(txData, this.chainId)

      tx.raw[6] = getBufferFromHex(this.chainId.toString(16))
      tx.raw[7] = Buffer.from([])
      tx.raw[8] = Buffer.from([])
      const result = await this._ledgerClientIfExists.signTransaction(
        fullDerivationPath,
        tx.serialize().toString("hex")
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
      if (chainIdFromV !== this.chainId && !isValidSignedV) {
        throw new Error("Invalid chainID")
      }

      await this._destroyLedgerClientAsync()
      return `0x${tx.serialize().toString("hex")}`
    } catch (error) {
      await this._destroyLedgerClientAsync()
      throw error
    }
  }

  async getAccountsAsync(numberOfAccounts, accountsOffSet = 0) {
    const addresses = []
    for (
      let index = accountsOffSet;
      index < numberOfAccounts + accountsOffSet;
      index++
    ) {
      const address = await this.getAddress(index)
      addresses.push(address)
    }

    return addresses
  }

  async getAddress(index) {
    const path = this._baseDerivationPath.replace("x", index)

    let ledgerResponse
    try {
      this._ledgerClientIfExists = await this._createLedgerClientAsync()
      ledgerResponse = await this._ledgerClientIfExists.getAddress(
        path,
        this._shouldAlwaysAskForConfirmation,
        true
      )
    } finally {
      await this._destroyLedgerClientAsync()
    }

    const address = web3Utils.toChecksumAddress(ledgerResponse.address)

    this.addressToPathMap[address] = path
    this.pathToAddressMap[path] = address

    return address
  }
}
