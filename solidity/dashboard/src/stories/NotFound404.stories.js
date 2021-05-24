import React from "react"
import centered from "@storybook/addon-centered/react"
import { NotFound404 } from "../components/NotFound404"

export default {
  title: "NotFound404",
  component: NotFound404,
  decorators: [centered],
}

const Template = (args) => <NotFound404 {...args} />

export const NoMatchForURL = Template.bind({})
NoMatchForURL.args = {}
