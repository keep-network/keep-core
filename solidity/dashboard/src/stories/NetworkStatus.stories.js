import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { NetworkStatus } from "../components/NetworkStatus"

storiesOf("NetworkStatus", module).addDecorator(centered)

export default {
  title: "NetworkStatus",
  component: NetworkStatus,
}

const Template = (args) => <NetworkStatus {...args} />

export const Default = Template.bind({})
Default.args = {}

// TODO: component with hook - do stories for network connected
