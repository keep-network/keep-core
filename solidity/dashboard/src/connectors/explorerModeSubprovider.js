import { Subprovider } from "@0x/subproviders"

class ExplorerModeSubprovider extends Subprovider {
  constructor(_eventEmitter) {
    super()
    this.eventEmitter = _eventEmitter
  }
  async handleRequest(payload, next, end) {
    switch (payload.method) {
      case "eth_sendTransaction":
        this.eventEmitter.emit("chooseWalletAndSendTransaction", payload)
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
