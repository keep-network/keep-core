import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import NoData from "../components/NoData"
import * as Icons from "../components/Icons"

storiesOf("NoData", module).addDecorator(centered)

export default {
  title: "NoData",
  component: NoData,
}

const Template = (args) => <NoData {...args} />

export const Default = Template.bind({})
Default.args = {
  title: "title",
  iconComponents: Icons.KeepBlackGreen,
  content: "content",
}
