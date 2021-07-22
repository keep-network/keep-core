import React from "react"
import centered from "@storybook/addon-centered/react"
import NavLink from "../components/NavLink"

export default {
  title: "NavLink",
  component: NavLink,
  argTypes: {
    onClick: {
      action: "onMessageClose function called",
    },
  },
  decorators: [centered],
}

const Template = (args) => <NavLink {...args} />

export const Default = Template.bind({})
Default.args = {
  to: "/liquidity",
  children: "link to subpage",
}
