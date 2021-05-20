import React from "react"
import { SpeechBubbleTooltip } from "../components/SpeechBubbleTooltip"
import { colors } from "../constants/colors"

export default {
  title: "SpeechBubbleTooltip",
  component: SpeechBubbleTooltip,
}

const Template = (args) => <SpeechBubbleTooltip {...args} />

export const Default = Template.bind({})
Default.args = {
  text: "text",
  title: "title",
  iconColor: colors.mint80,
  iconBackgroundColor: colors.mint20,
}

// TODO: Refactor SpeechBubbleTooltip component because it is not updated (add
//  Tooltip.Header!)
