import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import centered from "@storybook/addon-centered/react"
import { Rewards } from "../components/Rewards"

export default {
  title: "Rewards",
  component: Rewards,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <Rewards {...args} />

export const IsFetching = Template.bind({})
IsFetching.args = { }

// TODO: mock fetched data
