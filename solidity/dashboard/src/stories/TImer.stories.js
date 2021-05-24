import React from "react"
import Timer from "../components/Timer"
import centered from "@storybook/addon-centered/react"

export default {
  title: "Timer",
  component: Timer,
  decorators: [
    (Story) => (
      <section
        className="tile rewards-countdown"
        style={{ width: "50rem", height: "11rem" }}
      >
        <h2 className="h2--alt">
          Next rewards release:&nbsp;
          <Story />
        </h2>
      </section>
    ),
    centered,
  ],
}

const Template = (args) => <Timer {...args} />

export const RewardRelease = Template.bind({})
RewardRelease.args = { targetInUnix: "1620950400" }
