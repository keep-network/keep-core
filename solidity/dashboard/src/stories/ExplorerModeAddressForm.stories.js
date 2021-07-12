import React from "react"
import centered from "@storybook/addon-centered/react"
import ExplorerModeAddressForm from "../components/ExplorerModeAddressForm"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

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
  decorators: [whiteBackground, centered],
}

const Template = (args) => <ExplorerModeAddressForm {...args} />

export const Default = Template.bind({})
Default.args = { }