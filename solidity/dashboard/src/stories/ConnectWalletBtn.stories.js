import React from "react"
import centered from "@storybook/addon-centered/react"
import ConnectWalletBtn from "../components/ConnectWalletBtn"

export default {
  title: "ConnectWalletBtn",
  component: ConnectWalletBtn,
  decorators: [centered],
}

const Template = (args) => <ConnectWalletBtn {...args} />

export const Default = Template.bind({})
Default.args = {
  text: "connect wallet",
}

export const WithoutExplorerMode = Template.bind({})
WithoutExplorerMode.args = {
  text: "connect wallet",
  displayExplorerMode: false,
}
