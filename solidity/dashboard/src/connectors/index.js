import { TrezorProvider } from "./trezor"
import { LedgerProvider, LEDGER_DERIVATION_PATHS } from "./ledger"

const InjectedProvider = window.ethereum
if (InjectedProvider) {
  // https://docs.metamask.io/guide/ethereum-provider.html#ethereum-autorefreshonnetworkchange
  InjectedProvider.autoRefreshOnNetworkChange = false
}

export {
  TrezorProvider,
  LedgerProvider,
  LEDGER_DERIVATION_PATHS,
  InjectedProvider,
}
