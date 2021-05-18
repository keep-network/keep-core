import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import CurrentWalletIconTooltip from "../components/CurrentWalletIconTooltip"

// TODO: COMPONENT WITH HOOK

storiesOf("CurrentWalletIconTooltip", module).addDecorator(centered)

export default {
  title: "CurrentWalletIconTooltip",
  component: CurrentWalletIconTooltip,
}

const Template = (args) => <CurrentWalletIconTooltip {...args} />

export const Default = Template.bind({})
Default.args = { }
