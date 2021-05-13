import React from "react"
import Timer from "../components/Timer"

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
  ],
}

const Template = (args) => <Timer {...args} />

export const RewardRelease = Template.bind({})
RewardRelease.args = { targetInUnix: "1620950400" }
