import React from "react"
import centered from "@storybook/addon-centered/react"
import ConfirmationModal from "../components/ConfirmationModal"

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
  decorators: [centered],
}

const Template = (args) => <ConfirmationModal {...args} />

export const Default = Template.bind({})
Default.args = {
  confirmationText: "SUBMIT",
  btnText: "Submit",
}
