import React from "react"
import Footer from "../components/Footer"
import centered from "@storybook/addon-centered/react"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "Footer",
  component: Footer,
  decorators: [whiteBackground, centered],
}

const Template = (args) => <Footer {...args} />

export const Default = Template.bind({})
Default.args = { targetInUnix: "1620950400" }
