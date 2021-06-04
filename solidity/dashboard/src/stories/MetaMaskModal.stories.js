import React from "react"
import centered from "@storybook/addon-centered/react"
import { injected } from "../connectors"
import MetaMaskModal from "../components/MetaMaskModal"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

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
  decorators: [whiteBackground, centered],
}

const Template = (args) => <MetaMaskModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: injected,
}
