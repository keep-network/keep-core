import React from "react"
import centered from "@storybook/addon-centered/react"
import { NetworkStatus } from "../components/NetworkStatus"

export default {
  title: "NetworkStatus",
  component: NetworkStatus,
  decorators: [centered],
}

const Template = (args) => <NetworkStatus {...args} />

export const Default = Template.bind({})
Default.args = {}

// TODO: component with hook - do stories for network connected
