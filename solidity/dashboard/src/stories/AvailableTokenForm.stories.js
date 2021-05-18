import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import AvailableTokenForm from "../components/AvailableTokenForm"

// TODO: COMPONENT WITH HOOK

storiesOf("AvailableTokenForm", module).addDecorator(centered)

export default {
  title: "AvailableTokenForm",
  component: AvailableTokenForm,
  argTypes: {
    onSubmit: {
      action: "operator selected",
    },
  },
}

const Template = (args) => <AvailableTokenForm {...args} />

export const Default = Template.bind({})
Default.args = { submitBtnText: "submit text" }
