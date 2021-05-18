// TODO: Story with hook and Formik

import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import FormCheckbox from "../components/FormCheckbox"

storiesOf("FormCheckbox", module).addDecorator(centered)

export default {
  title: "FormCheckbox",
  component: FormCheckbox,
}

const Template = (args) => <FormCheckbox {...args} />

// export const Default = Template.bind({})
// Default.args = {
//   name: "checkboxTest",
//   type: "checkbox",
//   label: "Check this checkbox",
// }
