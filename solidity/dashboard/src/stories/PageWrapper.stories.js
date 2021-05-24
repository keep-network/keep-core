import React from "react"
import centered from "@storybook/addon-centered/react"
import PageWrapper from "../components/PageWrapper"

export default {
  title: "PageWrapper",
  component: PageWrapper,
  decorators: [centered],
}

const Template = (args) => <PageWrapper {...args} />

export const Default = Template.bind({})
Default.args = {}

export const WithTitle = Template.bind({})
WithTitle.args = { title: "PageWrapper title" }
