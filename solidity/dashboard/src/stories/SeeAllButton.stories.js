import React from "react"
import centered from "@storybook/addon-centered/react"
import { SeeAllButton } from "../components/SeeAllButton"

export default {
  title: "SeeAllButton",
  component: SeeAllButton,
  decorators: [centered],
}

const Template = (args) => <SeeAllButton {...args} />

export const SeeAll = Template.bind({})
SeeAll.args = {
  previewDataCount: 20,
  dataLength: 50,
  showAll: false,
}

export const SeeLess = Template.bind({})
SeeLess.args = {
  previewDataCount: 20,
  dataLength: 50,
  showAll: true,
}
