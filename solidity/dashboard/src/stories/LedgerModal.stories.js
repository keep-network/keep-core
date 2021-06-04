import React from "react"
import LedgerModal from "../components/LedgerModal"
import centered from "@storybook/addon-centered/react"
import { LEDGER_DERIVATION_PATHS, LedgerConnector } from "../connectors"
import { WALLETS } from "../constants/constants"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "LedgerModal",
  component: LedgerModal,
  argTypes: {
    connectAppWithWallet: {
      action: "connectAppWithWallet function called",
    },
    closeModal: {
      action: "closeModal clicked",
    },
  },
  decorators: [whiteBackground, centered],
}

const Template = (args) => <LedgerModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: {
    name: WALLETS.LEDGER.name,
    LEDGER_LIVE: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LIVE),
    LEDGER_LEGACY: new LedgerConnector(LEDGER_DERIVATION_PATHS.LEDGER_LEGACY),
  },
}
