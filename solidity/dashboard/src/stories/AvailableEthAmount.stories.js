import React from "react"
import AvailableEthAmount from "../components/AvailableEthAmount"
import centered from "@storybook/addon-centered/react"

export default {
  title: "AvailableEthAmount",
  component: AvailableEthAmount,
  decorators: [centered],
}

const Template = (args) => <AvailableEthAmount {...args} />

export const Default = Template.bind({})
Default.args = { availableETHInWei: "20000000000000000000000000" }
