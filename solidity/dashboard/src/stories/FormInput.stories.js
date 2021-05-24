// TODO: Story with hook and Formik

import React from "react"
import centered from "@storybook/addon-centered/react"
import FormInput from "../components/FormInput"

export default {
  title: "FormInput",
  component: FormInput,
  decorators: [centered],
}

const Template = (args) => <FormInput {...args} />

// export const Default = Template.bind({})
// Default.args = {
//   name: "amount",
//   type: "text",
//   label: "KEEP Amount",
//   placeholder: "0",
// }
