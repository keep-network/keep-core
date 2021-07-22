import React from "react"
import centered from "@storybook/addon-centered/react"
import AddressShortcut from "../components/AddressShortcut"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "AddressShortcut",
  component: AddressShortcut,
  argTypes: {
    onClick: {
      action: "onClick clicked",
    },
  },
  decorators: [whiteBackground, centered],
}

const Template = (args) => <AddressShortcut {...args} />

export const Default = Template.bind({})
Default.args = { address: "0x5777C7DdEd294654FbefC1Ed262fC8Ba4Ac40De1" }
