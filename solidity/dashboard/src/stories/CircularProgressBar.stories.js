import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import CircularProgressBar from "../components/CircularProgressBar"

storiesOf("CircularProgressBar", module).addDecorator(centered)

export default {
  title: "CircularProgressBar",
  component: CircularProgressBar,
}

const Template = (args) => <CircularProgressBar {...args} />

export const Default = Template.bind({})
Default.args = { value: 30, total: 100 }
