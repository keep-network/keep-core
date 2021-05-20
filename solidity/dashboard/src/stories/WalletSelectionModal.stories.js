import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { TrezorConnector } from "../connectors"
import TrezorModal from "../components/TrezorModal"
import { WalletSelectionModal } from "../components/WalletSelectionModal"

storiesOf("WalletSelectionModal", module).addDecorator(centered)

export default {
  title: "WalletSelectionModal",
  component: WalletSelectionModal,
}

const Template = (args) => <WalletSelectionModal {...args} />

export const Default = Template.bind({})
Default.args = {}
