import { Subprovider } from "@0x/subproviders";


class ExplorerModeSubprovider extends Subprovider {
  constructor(_eventEmitter) {
    super()
    this.eventEmitter = _eventEmitter
  }
  async handleRequest(payload, next, end) {
    switch (payload.method) {
      case "eth_sendTransaction":
        this.eventEmitter.emit("chooseWalletAndSendTransaction", payload)
        end()
        return
      default:
        next()
        return
    }
  }
}

export default ExplorerModeSubprovider
