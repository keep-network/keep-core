import React from "react"
import store from "../store"
import { Provider } from "react-redux"
import ResourceTooltip from "../components/ResourceTooltip"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

// storiesOf("ResourceTooltip", module).addDecorator(centered)

export default {
  title: "ResourceTooltip",
  component: ResourceTooltip,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
  ],
}

const Template = (args) => <ResourceTooltip {...args} />

export const Default = Template.bind({})
Default.args = {
  title: "Title",
  content: "content",
  withRedirectButton: false,
}
