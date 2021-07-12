import React from "react"
import centered from "@storybook/addon-centered/react"
import Loadable from "../components/Loadable"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "Loadable",
  component: Loadable,
  decorators: [whiteBackground, centered],
}

const Template = (args) => <Loadable {...args} />

export const IsFetching = Template.bind({})
IsFetching.args = {
  isFetching: false,
  children: "text",
}
