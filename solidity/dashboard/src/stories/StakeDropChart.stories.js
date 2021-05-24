import React from "react"
import StakeDropChart from "../components/StakeDropChart"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react";

/**
 * StakeDropChart is dropped for now, so we are not displaying story for it
 */

storiesOf("StakeDropChart", module).addDecorator(centered)

export default {
  title: "StakeDropChart",
  component: StakeDropChart,
  decorators: [
    (Story) => (
      <section className="rewards-overview--random-beacon">
        <section>
          <Story />
        </section>
      </section>
    ),
  ],
}

// const Template = (args) => <StakeDropChart {...args} />

// export const Default = Template.bind({})
// Default.args = {}
