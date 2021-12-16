import { WALLETS } from "../constants/constants"
import { InjectedConnector } from "./injected-connector"

class TallyInjectedConnector extends InjectedConnector {
  constructor() {
    super(WALLETS.TALLY.name, "tally")
  }
}

export { TallyInjectedConnector }
