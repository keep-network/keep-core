import React from "react"
import Card from "../components/Card"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("Card", module).addDecorator(centered)

export default {
  title: "Card",
  component: Card,
}

const Template = (args) => <Card {...args} />

export const Default = Template.bind({})
Default.args = { children: "Card content", className: "tile" }
