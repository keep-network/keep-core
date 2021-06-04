import React from "react"
import centered from "@storybook/addon-centered/react"
import { TrezorConnector } from "../connectors"
import TrezorModal from "../components/TrezorModal"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "TrezorModal",
  component: TrezorModal,
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

const Template = (args) => <TrezorModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: new TrezorConnector(),
}
