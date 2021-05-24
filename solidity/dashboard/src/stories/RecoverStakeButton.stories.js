import React from "react"
import RecoverStakeButton from "../components/RecoverStakeButton"
import store from "../store"
import { Provider } from "react-redux"
import centered from "@storybook/addon-centered/react"

export default {
  title: "RecoverStakeButton",
  component: RecoverStakeButton,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    centered,
  ],
}

const Template = (args) => <RecoverStakeButton {...args} />

export const Default = Template.bind({})
Default.args = {
  operatrorAddress: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
  btnClassName: "btn btn-sm btn-secondary",
  btnText: "recover",
}
