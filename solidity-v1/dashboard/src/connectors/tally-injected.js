import { WALLETS } from "../constants/constants"
import { InjectedConnector } from "./injected-connector"

class TallyInjectedConnector extends InjectedConnector {
  constructor() {
    super(WALLETS.TALLY.name, window.tally)
  }
}

export { TallyInjectedConnector }
