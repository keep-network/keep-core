import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import ConnectWalletBtn from "../components/ConnectWalletBtn"

storiesOf("ConnectWalletBtn", module).addDecorator(centered)

export default {
  title: "ConnectWalletBtn",
  component: ConnectWalletBtn,
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
