import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import PageWrapper from "../components/PageWrapper"

storiesOf("PageWrapper", module).addDecorator(centered)

export default {
  title: "PageWrapper",
  component: PageWrapper,
}

const Template = (args) => <PageWrapper {...args} />

export const Default = Template.bind({})
Default.args = {}

export const WithTitle = Template.bind({})
WithTitle.args = { title: "PageWrapper title" }
