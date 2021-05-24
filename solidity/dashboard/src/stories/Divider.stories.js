import React from "react"
import centered from "@storybook/addon-centered/react"
import Divider from "../components/Divider"

export default {
  title: "Divider",
  component: Divider,
  decorators: [centered],
}

const Template = (args) => <Divider {...args} />

export const GreyDivider = Template.bind({})
GreyDivider.args = { style: { borderTop: "1px solid grey", height: "30px" } }
