import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import { ECDSARewardsDetails } from "../components/RewardsDetails"
import centered from "@storybook/addon-centered/react";

export default {
  title: "ECDSARewardsDetails",
  component: ECDSARewardsDetails,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <section className={"tile"} style={{ width: "20rem" }}>
          <Story />
        </section>
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <ECDSARewardsDetails {...args} />

export const Default = Template.bind({})
Default.args = {}