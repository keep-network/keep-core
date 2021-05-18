import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import Loadable from "../components/Loadable"

storiesOf("Loadable", module).addDecorator(centered)

export default {
  title: "Loadable",
  component: Loadable,
}

const Template = (args) => <Loadable {...args} />

export const IsFetching = Template.bind({})
IsFetching.args = {
  isFetching: false,
  children: "text",
}
