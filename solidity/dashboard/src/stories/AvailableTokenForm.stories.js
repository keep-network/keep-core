import React from "react"
import centered from "@storybook/addon-centered/react"
import AvailableTokenForm from "../components/AvailableTokenForm"
import withFormik from "storybook-formik"
import {
  formatAmount as formatFormAmount,
  normalizeAmount,
} from "../forms/form.utils"
import MaxAmountAddon from "../components/MaxAmountAddon"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

// TODO: Get rid of error <form> cannot appear as descendant of <form>

export default {
  title: "AvailableTokenForm",
  component: AvailableTokenForm,
  argTypes: {
    onSubmit: {
      action: "operator selected",
    },
  },
  decorators: [whiteBackground, centered, withFormik],
}

const Template = (args) => <AvailableTokenForm {...args} />

export const Default = Template.bind({})
Default.args = {
  submitBtnText: "submit text",
  formInputProps: {
    name: "amount",
    type: "text",
    label: "Withdraw",
    normalize: normalizeAmount,
    format: formatFormAmount,
    placeholder: "0",
    inputAddon: <MaxAmountAddon text="Max KEEP" />,
  },
}
