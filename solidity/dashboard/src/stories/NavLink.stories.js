import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import NavLink from "../components/NavLink"

storiesOf("NavLink", module).addDecorator(centered)

export default {
  title: "NavLink",
  component: NavLink,
  argTypes: {
    onClick: {
      action: "onMessageClose function called",
    },
  },
}

const Template = (args) => <NavLink {...args} />

export const Default = Template.bind({})
Default.args = {
  to: "/liquidity",
  children: "link to subpage",
}
