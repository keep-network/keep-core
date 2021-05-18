import React from "react"
import Button from "../components/Button"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("Button", module).addDecorator(centered)

export default {
  title: "Button",
  component: Button,
  decorators: [
    (Story) => (
      <div style={{ height: "20px" }}>
        <Story />
      </div>
    ),
  ],
}

const Template = (args) => <Button {...args} />

export const Primary = Template.bind({})
Primary.args = { children: "Click me!", className: "btn btn-lg btn-primary" }

export const Secondary = Template.bind({})
Secondary.args = {
  children: "Click me!",
  className: "btn btn-lg btn-secondary",
}

export const PrimaryIsFetching = Template.bind({})
PrimaryIsFetching.args = {
  children: "Click me!",
  className: "btn btn-lg btn-primary",
  isFetching: true,
}

export const SecondaryIsFetching = Template.bind({})
SecondaryIsFetching.args = {
  children: "Click me!",
  className: "btn btn-lg btn-secondary",
  isFetching: true,
}
