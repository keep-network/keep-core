import React from "react"
import PageWrapper from "../components/PageWrapper"

export default {
  title: "PageWrapper",
  component: PageWrapper,
}

const Template = (args) => <PageWrapper {...args} />

export const Default = Template.bind({})
Default.args = {}

export const WithTitle = Template.bind({})
WithTitle.args = { title: "PageWrapper title" }
