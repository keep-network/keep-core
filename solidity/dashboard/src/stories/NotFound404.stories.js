import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import * as Icons from "../components/Icons"
import { NotFound404 } from "../components/NotFound404"

storiesOf("NotFound404", module).addDecorator(centered)

export default {
  title: "NotFound404",
  component: NotFound404,
}

const Template = (args) => <NotFound404 {...args} />

export const NoMatchForURL = Template.bind({})
NoMatchForURL.args = {}
