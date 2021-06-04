import React from "react"
import centered from "@storybook/addon-centered/react"
import { WalletConnectConnector } from "../connectors"
import WalletConnectModal from "../components/WalletConnectModal"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "WalletConnectModal",
  component: WalletConnectModal,
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

const Template = (args) => <WalletConnectModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: new WalletConnectConnector(),
}
