import React from "react"
import Card from "../components/Card"
import centered from "@storybook/addon-centered/react"

export default {
  title: "Card",
  component: Card,
  decorators: [centered],
}

const Template = (args) => <Card {...args} />

export const Default = Template.bind({})
Default.args = { children: "Card content", className: "tile" }
