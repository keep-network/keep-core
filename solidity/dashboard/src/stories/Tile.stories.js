import React from "react"
import Tile from "../components/Tile"

export default {
  title: "Tile",
  component: Tile,
  decorators: [
    (Story) => (
      <div style={{ width: "30rem" }}>
        <Story />
      </div>
    ),
  ],
}

const Template = (args) => <Tile {...args} />

export const Default = Template.bind({})
Default.args = {
  title: "title",
  subtitle: "subtitle",
  children: "content",
  withTooltip: false,
  tooltipProps: {
    text: "",
    title: "",
  },
}

export const WithTooltip = Template.bind({})
WithTooltip.args = {
  title: "title",
  subtitle: "subtitle",
  children: "content",
  withTooltip: true,
  tooltipProps: {
    text: "tooltip text",
    title: "Tooltip Title",
  },
}
