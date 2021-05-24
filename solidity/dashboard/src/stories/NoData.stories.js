import React from "react"
import centered from "@storybook/addon-centered/react"
import NoData from "../components/NoData"
import * as Icons from "../components/Icons"

export default {
  title: "NoData",
  component: NoData,
  decorators: [centered],
}

const Template = (args) => <NoData {...args} />

export const Default = Template.bind({})
Default.args = {
  title: "title",
  iconComponents: Icons.KeepBlackGreen,
  content: "content",
}
