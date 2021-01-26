import { TrezorProvider } from "./trezor"
import { LedgerProvider, LEDGER_DERIVATION_PATHS } from "./ledger"

const InjectedProvider = window.ethereum
if (InjectedProvider) {
  // https://docs.metamask.io/guide/ethereum-provider.html#ethereum-autorefreshonnetworkchange
  InjectedProvider.autoRefreshOnNetworkChange = false

  // Override `enable` function in the injected provider.
  // MetaMask: 'ethereum.enable()' is deprecated and may be removed in the
  // future. https://docs.metamask.io/guide/provider-migration.html#replacing-window-web3
  InjectedProvider.enable = async () => {
    return await InjectedProvider.request({ method: "eth_requestAccounts" })
  }
}

export {
  TrezorProvider,
  LedgerProvider,
  LEDGER_DERIVATION_PATHS,
  InjectedProvider,
}
