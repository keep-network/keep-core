import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import { BeaconRewardsDetails } from "../components/RewardsDetails"
import centered from "@storybook/addon-centered/react";

export default {
  title: "BeaconRewardsDetails",
  component: BeaconRewardsDetails,
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

const Template = (args) => <BeaconRewardsDetails {...args} />

export const Default = Template.bind({})
Default.args = {}
