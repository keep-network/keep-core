import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import ConfirmationModal from "../components/ConfirmationModal"

storiesOf("ConfirmationModal", module).addDecorator(centered)

export default {
  title: "ConfirmationModal",
  component: ConfirmationModal,
  argTypes: {
    onBtnClick: {
      action: "onBtnClick clicked",
    },
    onCancel: {
      action: "onCancel clicked",
    },
  },
}

const Template = (args) => <ConfirmationModal {...args} />

export const Default = Template.bind({})
Default.args = {
  confirmationText: "SUBMIT",
  btnText: "Submit",
}
