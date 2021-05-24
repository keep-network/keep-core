import React from "react"
import { SideMenu } from "../components/SideMenu"
import centered from "@storybook/addon-centered/react"

export default {
  title: "SideMenu",
  component: SideMenu,
  decorators: [
    (Story) => (
      <div style={{ width: "20rem" }}>
        <Story />
      </div>
    ),
    centered,
  ],
}

const Template = (args) => <SideMenu {...args} />

export const Default = Template.bind({})
Default.args = {}
