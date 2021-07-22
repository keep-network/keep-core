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
GreyDivider.args = {
  style: {
    borderTop: "5px solid grey",
    margin: "2rem -2rem 0",
    padding: "2rem 2rem 0",
  },
}
