import { TrezorConnector } from "./trezor"
import { LedgerConnector, LEDGER_DERIVATION_PATHS } from "./ledger"
import { InjectedConnector } from "./injected"
import { WalletConnectConnector } from "./wallet-connect"
import { UserRejectedConnectionRequestError } from "./utils"

const injected = new InjectedConnector()

export {
  TrezorConnector,
  LedgerConnector,
  LEDGER_DERIVATION_PATHS,
  injected,
  WalletConnectConnector,
  UserRejectedConnectionRequestError,
}
