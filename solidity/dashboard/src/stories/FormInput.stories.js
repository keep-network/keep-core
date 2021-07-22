import React from "react"
import centered from "@storybook/addon-centered/react"
import FormInput from "../components/FormInput"
import withFormik from "storybook-formik"

export default {
  title: "FormInput",
  component: FormInput,
  decorators: [centered, withFormik],
}

const Template = (args) => <FormInput {...args} />

export const Default = Template.bind({})
Default.args = {
  name: "name",
  type: "text",
  label: "label",
  placeholder: "placeholder",
}
