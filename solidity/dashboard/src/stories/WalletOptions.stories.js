import React from "react"
import centered from "@storybook/addon-centered/react"
import WalletOptions from "../components/WalletOptions"
import { blackBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "WalletOptions",
  component: WalletOptions,
  decorators: [blackBackground, centered],
}

const Template = (args) => <WalletOptions {...args} />

export const Default = Template.bind({})
Default.args = {
  displayExplorerMode: true,
}

export const WithoutExplorerMode = Template.bind({})
WithoutExplorerMode.args = {
  displayExplorerMode: false,
}
