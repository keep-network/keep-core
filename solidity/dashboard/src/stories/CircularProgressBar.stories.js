import React from "react"
import centered from "@storybook/addon-centered/react"
import CircularProgressBar from "../components/CircularProgressBar"

export default {
  title: "CircularProgressBar",
  component: CircularProgressBar,
  decorators: [centered],
}

const Template = (args) => <CircularProgressBar {...args} />

export const Default = Template.bind({})
Default.args = { value: 30, total: 100 }
