import React from "react"
import centered from "@storybook/addon-centered/react"
import { Web3Status } from "../components/Web3Status"

export default {
  title: "Web3Status",
  component: Web3Status,
  decorators: [centered]
}

const Template = (args) => <Web3Status {...args} />

export const NotConnected = Template.bind({})
NotConnected.args = {
  yourAddress: null,
  isConnected: false,
  connector: null,
}

// TODO: Connected - component with hook (useWeb3Context)
// export const Connected = Template.bind({})
// Connected.args = {
//   yourAddress: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
//   isConnected: true,
//   connector: injected,
// }
