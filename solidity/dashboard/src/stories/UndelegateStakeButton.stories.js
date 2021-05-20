import React from "react"
import UndelegateStakeButton from "../components/UndelegateStakeButton"
import { Provider } from "react-redux"
import store from "../store"

export default {
  title: "UndelegateStakeButton",
  component: UndelegateStakeButton,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
  ],
}

const Template = (args) => <UndelegateStakeButton {...args} />

export const Default = Template.bind({})
Default.args = {
  operator: "0xd7E826Ae811942142FBe350d68b6171937Ac408f",
  isFromGrant: false,
  isInInitializationPeriod: false,
}
