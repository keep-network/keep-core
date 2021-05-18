// TODO: Story with hook and Formik

import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import FormInput from "../components/FormInput"

storiesOf("FormInput", module).addDecorator(centered)

export default {
  title: "FormInput",
  component: FormInput,
}

const Template = (args) => <FormInput {...args} />

// export const Default = Template.bind({})
// Default.args = {
//   name: "amount",
//   type: "text",
//   label: "KEEP Amount",
//   placeholder: "0",
// }
