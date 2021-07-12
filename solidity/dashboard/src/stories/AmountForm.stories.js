import React from "react"
import centered from "@storybook/addon-centered/react"
import AmountForm from "../components/AmountForm"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "AmountForm",
  component: AmountForm,
  argTypes: {
    onCancel: {
      action: "onCancel clicked",
    },
    onBtnClick: {
      action: "onBtnClick clicked",
    },
  },
  decorators: [whiteBackground, centered],
}

const Template = (args) => <AmountForm {...args} />

export const Default = Template.bind({})
Default.args = {
  submitBtnText: "add keep",
  availableAmount: "300000000000000000000",
  currentAmount: "300000000000000000000",
}
