import React from "react"
import centered from "@storybook/addon-centered/react"
import Badge from "../components/Badge"

export default {
  title: "Badge",
  component: Badge,
  decorators: [centered],
}

const Template = (args) => <Badge {...args} />

export const Primary = Template.bind({})
Primary.args = { text: "badge" }
