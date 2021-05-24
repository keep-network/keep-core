import React from "react"
import centered from "@storybook/addon-centered/react"
import { Web3StatusView } from "../components/Web3Status"
import { injected } from "../connectors"
import { Provider } from "react-redux"
import store from "../store"

export default {
  title: "Web3StatusView",
  component: Web3StatusView,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <Web3StatusView {...args} />

export const NotConnected = Template.bind({})
NotConnected.args = {
  yourAddress: null,
  isConnected: false,
  connector: null,
}

export const Connected = Template.bind({})
Connected.args = {
  yourAddress: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
  isConnected: true,
  connector: injected,
}
