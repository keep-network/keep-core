import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import WalletOptions from "../components/WalletOptions"

storiesOf("WalletConnectModal", module).addDecorator(centered)

export default {
  title: "WalletOptions",
  component: WalletOptions,
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
