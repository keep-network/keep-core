import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import DelegatedTokens from "../components/DelegatedTokens"
import { Provider } from "react-redux"
import store from "../store"

const mockData = {
  stakedBalance: "20000000000000000000000000",
  ownerAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
  beneficiaryAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
  authorizerAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
  error: null,
  slashedTokensError: null,
  isDelegationFromGrant: true,
  isInInitializationPeriod: false,
  undelegationPeriod: "1209600",
  isManagedGrant: false,
  undelegationCompletedAt: null,
}

storiesOf("DelegatedTokens", module).addDecorator(centered)

export default {
  title: "DelegatedTokens",
  component: DelegatedTokens,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
  ],
  argTypes: {
    onSubmit: {
      action: "onSubmit clicked",
    },
  },
}

const Template = (args) => <DelegatedTokens {...args} />

export const WithMockedData = Template.bind({})
WithMockedData.args = { data: mockData }
