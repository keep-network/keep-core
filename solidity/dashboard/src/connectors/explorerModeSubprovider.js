import { Subprovider } from "@0x/subproviders"
import { EventEmitter } from "events"

class ExplorerModeSubprovider extends Subprovider {
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

Object.assign(ExplorerModeSubprovider.prototype, EventEmitter.prototype)

export default ExplorerModeSubprovider
