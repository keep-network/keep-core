import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { WalletConnectConnector } from "../connectors"
import WalletConnectModal from "../components/WalletConnectModal"

storiesOf("WalletConnectModal", module).addDecorator(centered)

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
}

const Template = (args) => <WalletConnectModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: new WalletConnectConnector(),
}
