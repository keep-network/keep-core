// TODO: Get rid of error 'A component is changing an uncontrolled input of type checkbox to be controlled'

import React from "react"
import centered from "@storybook/addon-centered/react"
import FormCheckbox from "../components/FormCheckbox"
import withFormik from "storybook-formik"

export default {
  title: "FormCheckbox",
  component: FormCheckbox,
  decorators: [centered, withFormik],
}

const Template = (args) => <FormCheckbox {...args} />

export const Default = Template.bind({})
Default.args = {
  name: "checkboxTest",
  type: "checkbox",
  label: "Check this checkbox",
}
