import React from "react"
import centered from "@storybook/addon-centered/react"
import { WalletSelectionModal } from "../components/WalletSelectionModal"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "WalletSelectionModal",
  component: WalletSelectionModal,
  decorators: [whiteBackground, centered],
}

const Template = (args) => <WalletSelectionModal {...args} />

export const Default = Template.bind({})
Default.args = {}
