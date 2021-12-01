import { TrezorConnector } from "./trezor"
import { LedgerConnector, LEDGER_DERIVATION_PATHS } from "./ledger"
import { MetaMaskInjectedConnector } from "./metamask-injected"
import { WalletConnectConnector } from "./wallet-connect"
import { UserRejectedConnectionRequestError } from "./utils"

const metaMaskInjectedConnector = new MetaMaskInjectedConnector()

export {
  TrezorConnector,
  LedgerConnector,
  LEDGER_DERIVATION_PATHS,
  metaMaskInjectedConnector,
  WalletConnectConnector,
  UserRejectedConnectionRequestError,
}
