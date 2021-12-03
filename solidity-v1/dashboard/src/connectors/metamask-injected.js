import { WALLETS } from "../constants/constants"
import { InjectedConnector } from "./injected-connector"

class MetaMaskInjectedConnector extends InjectedConnector {
  constructor() {
    super(WALLETS.METAMASK.name, window.ethereum)
  }
}

export { MetaMaskInjectedConnector }
