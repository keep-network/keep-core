import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import ResourceTooltip from "../components/ResourceTooltip"
import centered from "@storybook/addon-centered/react"

export default {
  title: "ResourceTooltip",
  component: ResourceTooltip,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <ResourceTooltip {...args} />

export const Default = Template.bind({})
Default.args = {
  title: "Title",
  content: "content",
  withRedirectButton: false,
}
