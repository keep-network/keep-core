import React from "react"
import centered from "@storybook/addon-centered/react"
import { WalletSelectionModal } from "../components/WalletSelectionModal"

export default {
  title: "WalletSelectionModal",
  component: WalletSelectionModal,
  decorators: [centered],
}

const Template = (args) => <WalletSelectionModal {...args} />

export const Default = Template.bind({})
Default.args = {}
