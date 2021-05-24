import React from "react"
import Tag from "../components/Tag"
import centered from "@storybook/addon-centered/react"
import * as Icons from "../components/Icons"

export default {
  title: "Tag",
  component: Tag,
  decorators: [centered],
}

const Template = (args) => <Tag {...args} />

export const KeepCurrent = Template.bind({})
KeepCurrent.args = { text: "Current", IconComponent: Icons.KeepToken }

export const Issued = Template.bind({})
Issued.args = { IconComponent: Icons.Time, text: "Issued" }
