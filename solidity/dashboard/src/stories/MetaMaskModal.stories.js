import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { injected } from "../connectors"
import MetaMaskModal from "../components/MetaMaskModal"

storiesOf("MetaMaskModal", module).addDecorator(centered)

export default {
  title: "MetaMaskModal",
  component: MetaMaskModal,
  argTypes: {
    connectAppWithWallet: {
      action: "connectAppWithWallet function called",
    },
    closeModal: {
      action: "closeModal clicked",
    },
  },
}

const Template = (args) => <MetaMaskModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: injected,
}
