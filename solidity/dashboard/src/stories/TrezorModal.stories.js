import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { TrezorConnector } from "../connectors"
import TrezorModal from "../components/TrezorModal";

storiesOf("TrezorModal", module).addDecorator(centered)

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
}

const Template = (args) => <TrezorModal {...args} />

export const Default = Template.bind({})
Default.args = {
  connector: new TrezorConnector(),
}
