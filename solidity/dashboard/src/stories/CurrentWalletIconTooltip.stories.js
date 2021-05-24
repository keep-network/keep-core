import React from "react"
import centered from "@storybook/addon-centered/react"
import CurrentWalletIconTooltip from "../components/CurrentWalletIconTooltip"

// TODO: COMPONENT WITH HOOK

export default {
  title: "CurrentWalletIconTooltip",
  component: CurrentWalletIconTooltip,
  decorators: [centered],
}

const Template = (args) => <CurrentWalletIconTooltip {...args} />

export const Default = Template.bind({})
Default.args = { }
