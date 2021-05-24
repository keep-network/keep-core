import React from "react"
import centered from "@storybook/addon-centered/react"
import { CurrentWalletIconTooltipView } from "../components/CurrentWalletIconTooltip"
import {
  injected,
  LEDGER_DERIVATION_PATHS,
  LedgerConnector,
  TrezorConnector,
  WalletConnectConnector,
} from "../connectors"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import { WALLETS } from "../constants/constants"

// TODO: MOVE WALLETS_OPTIONS TO A SEPARATE FILE AND EXPORT IT FOR EACH WALLET

export default {
  title: "CurrentWalletIconTooltipView",
  component: CurrentWalletIconTooltipView,
  decorators: [centered],
}

const Template = (args) => <CurrentWalletIconTooltipView {...args} />

export const MetaMask = Template.bind({})
MetaMask.args = { connector: injected }

export const Ledger = Template.bind({})
Ledger.args = {
  connector: {
    name: WALLETS.LEDGER.name,
    LEDGER_LIVE: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
    LEDGER_LEGACY: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
  },
}

export const Trezor = Template.bind({})
Trezor.args = { connector: new TrezorConnector() }

export const WalletConnect = Template.bind({})
WalletConnect.args = { connector: new WalletConnectConnector() }

export const ExplorerMode = Template.bind({})
ExplorerMode.args = { connector: new ExplorerModeConnector() }
