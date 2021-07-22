import React from "react"
import Tooltip from "../components/Tooltip"
import * as Icons from "../components/Icons"
import { colors } from "../constants/colors"
import centered from "@storybook/addon-centered/react"

export default {
  title: "Tooltip",
  component: Tooltip,
  decorators: [centered],
}

const Template = (args) => <Tooltip {...args} />

export const IconTooltip = Template.bind({})
IconTooltip.args = {
  children: "tooltip content",
  direction: "bottom",
  simple: false,
  triggerComponent: () => (
    <Icons.Tooltip color={colors.mint80} backgroundColor={colors.mint20} />
  ),
}

export const TextTooltip = Template.bind({})
TextTooltip.args = {
  children: "tooltip content",
  direction: "bottom",
  simple: false,
  triggerComponent: () => <div>Hover me!</div>,
}
