import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import SelectedWalletModal from "../components/SelectedWalletModal"
import * as Icons from "../components/Icons"
import { injected } from "../connectors"

storiesOf("SelectedWalletModal", module).addDecorator(centered)

export default {
  title: "SelectedWalletModal",
  component: SelectedWalletModal,
  argTypes: {
    connectAppWithWallet: {
      action: "connectAppWithWallet function called",
    },
    closeModal: {
      action: "closeModal clicked",
    },
  },
}

const Template = (args) => <SelectedWalletModal {...args} />

export const Default = Template.bind({})
Default.args = {
  icon: (
    <Icons.Diamond
      className="wallet-connect-logo wallet-connect-logo--black"
      width={30}
      height={28}
    />
  ),
  walletName: "WALLET NAME",
  connector: injected,
  connectWithWalletOnMount: true,
}
