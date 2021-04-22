/**
 * @dev ExplorerModeSubprovider class should have methods from EventEmitter
 * (from 'events' package) and Subprovider (from '@0x/subproviders' package)
 *
 * Since multimple inheritance in js is not possible, to make sure
 * ExplorerModeSubprovider works properly you should assign methods and
 * properties from both class mentioned above when initializing
 *
 * e.g:
 * const explorerModeSubprovider = Object.assign(
 *  ExplorerModeSubprovider.prototype,
 *  Subprovider.prototype,
 *  EventEmitter.prototype
 * )
 */
class ExplorerModeSubprovider {
  async handleRequest(payload, next, end) {
    switch (payload.method) {
      case "eth_sendTransaction":
        this.emit("chooseWalletAndSendTransaction", payload)
        const error = new Error(
          "You have to be connected to a wallet to do the transaction"
        )
        error.name = "ExplorerModeSubproviderError"
        end(error)
        return
      default:
        next()
        return
    }
  }
}

export default ExplorerModeSubprovider
