// TODO: Story with hook and Formik

import React from "react"
import centered from "@storybook/addon-centered/react"
import FormCheckbox from "../components/FormCheckbox"

export default {
  title: "FormCheckbox",
  component: FormCheckbox,
  decorators: [centered],
}

const Template = (args) => <FormCheckbox {...args} />

// export const Default = Template.bind({})
// Default.args = {
//   name: "checkboxTest",
//   type: "checkbox",
//   label: "Check this checkbox",
// }
