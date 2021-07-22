import React from "react"
import AvailableEthAmount from "../components/AvailableEthAmount"
import centered from "@storybook/addon-centered/react"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "AvailableEthAmount",
  component: AvailableEthAmount,
  decorators: [whiteBackground, centered],
}

const Template = (args) => <AvailableEthAmount {...args} />

export const Default = Template.bind({})
Default.args = { availableETHInWei: "20000000000000000000000000" }
