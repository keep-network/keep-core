import React from "react"
import centered from "@storybook/addon-centered/react"
import { NetworkStatusView } from "../components/NetworkStatus"

export default {
  title: "NetworkStatusView",
  component: NetworkStatusView,
  decorators: [centered],
}

const Template = (args) => <NetworkStatusView {...args} />

export const Connected = Template.bind({})
Connected.args = {
  networkType: "Network type",
  error: "",
  isConnected: true,
}

export const ConnectingError = Template.bind({})
ConnectingError.args = {
  networkType: "Network type",
  error: "This is error",
  isConnected: false,
}
