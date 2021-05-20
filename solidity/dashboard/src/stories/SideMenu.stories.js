import React from "react"
import { SideMenu } from "../components/SideMenu"

export default {
  title: "SideMenu",
  component: SideMenu,
  decorators: [
    (Story) => (
      <div style={{ width: "20rem" }}>
        <Story />
      </div>
    ),
  ],
}

const Template = (args) => <SideMenu {...args} />

export const Default = Template.bind({})
Default.args = {}
