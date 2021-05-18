import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import Badge from "../components/Badge"

storiesOf("Badge", module).addDecorator(centered)

export default {
  title: "Badge",
  component: Badge,
}

const Template = (args) => <Badge {...args} />

export const Primary = Template.bind({})
Primary.args = { text: "badge" }
