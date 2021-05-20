import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { Rewards } from "../components/Rewards"

storiesOf("Rewards", module).addDecorator(centered)

export default {
  title: "Rewards",
  component: Rewards,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
  ],
}

const Template = (args) => <Rewards {...args} />

export const IsFetching = Template.bind({})
IsFetching.args = { }

// TODO: mock fetched data
