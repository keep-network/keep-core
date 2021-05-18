import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import ExplorerModeAddressForm from "../components/ExplorerModeAddressForm"

storiesOf("ExplorerModeAddressForm", module).addDecorator(centered)

export default {
  title: "ExplorerModeAddressForm",
  component: ExplorerModeAddressForm,
  argTypes: {
    submitAction: {
      action: "submitAction",
    },
    onCancel: {
      action: "onCancel clicked",
    },
  },
}

const Template = (args) => <ExplorerModeAddressForm {...args} />

export const Default = Template.bind({})
Default.args = { }