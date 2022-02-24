import { TrezorConnector } from "./trezor"
import { LedgerConnector, LEDGER_DERIVATION_PATHS } from "./ledger"
import { MetaMaskInjectedConnector } from "./metamask-injected"
import { WalletConnectConnector } from "./wallet-connect"
import { UserRejectedConnectionRequestError } from "./utils"
import { TallyInjectedConnector } from "./tally-injected"

const metaMaskInjectedConnector = new MetaMaskInjectedConnector()
const tallyInjectedConnector = new TallyInjectedConnector()

export {
  TrezorConnector,
  LedgerConnector,
  LEDGER_DERIVATION_PATHS,
  metaMaskInjectedConnector,
  tallyInjectedConnector,
  WalletConnectConnector,
  UserRejectedConnectionRequestError,
}
