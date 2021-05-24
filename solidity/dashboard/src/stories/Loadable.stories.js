import React from "react"
import centered from "@storybook/addon-centered/react"
import Loadable from "../components/Loadable"

export default {
  title: "Loadable",
  component: Loadable,
  decorators: [centered],
}

const Template = (args) => <Loadable {...args} />

export const IsFetching = Template.bind({})
IsFetching.args = {
  isFetching: false,
  children: "text",
}
