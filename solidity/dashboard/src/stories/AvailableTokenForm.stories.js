import React from "react"
import centered from "@storybook/addon-centered/react"
import AvailableTokenForm from "../components/AvailableTokenForm"

// TODO: COMPONENT WITH HOOK

export default {
  title: "AvailableTokenForm",
  component: AvailableTokenForm,
  argTypes: {
    onSubmit: {
      action: "operator selected",
    },
  },
  decorators: [centered],
}

const Template = (args) => <AvailableTokenForm {...args} />

export const Default = Template.bind({})
Default.args = { submitBtnText: "submit text" }
