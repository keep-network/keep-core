import React from "react"
import AvailableEthAmount from "../components/AvailableEthAmount"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("AvailableEthAmount", module).addDecorator(centered)

export default {
  title: "AvailableEthAmount",
  component: AvailableEthAmount,
}

const Template = (args) => <AvailableEthAmount {...args} />

export const Default = Template.bind({})
Default.args = { availableETHInWei: "20000000000000000000000000" }
