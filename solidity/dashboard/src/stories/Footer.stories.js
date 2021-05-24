import React from "react"
import Footer from "../components/Footer"
import centered from "@storybook/addon-centered/react"

export default {
  title: "Footer",
  component: Footer,
  decorators: [centered],
}

const Template = (args) => <Footer {...args} />

export const Default = Template.bind({})
Default.args = { targetInUnix: "1620950400" }
