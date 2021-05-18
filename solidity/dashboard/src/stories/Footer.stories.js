import React from "react"
import Footer from "../components/Footer"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("Footer", module).addDecorator(centered)

export default {
  title: "Footer",
  component: Footer,
}

const Template = (args) => <Footer {...args} />

export const Default = Template.bind({})
Default.args = { targetInUnix: "1620950400" }
