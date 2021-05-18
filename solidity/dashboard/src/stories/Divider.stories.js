import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import Divider from "../components/Divider";

storiesOf("Divider", module).addDecorator(centered)

export default {
  title: "Divider",
  component: Divider,
}

const Template = (args) => <Divider {...args} />

export const GreyDivider = Template.bind({})
GreyDivider.args = { style: { borderTop: "1px solid grey", height: "30px" } }
