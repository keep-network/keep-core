import React from "react"
import centered from "@storybook/addon-centered/react"
import WithdrawETHModal from "../components/WithdrawETHModal"
import { Provider } from "react-redux"
import store from "../store"
import { whiteBackground } from "../../.storybook/cuatomDecorators"

export default {
  title: "WithdrawETHModal",
  component: WithdrawETHModal,
  decorators: [
    (Story) => (
      <Provider store={store}>
        <Story />
      </Provider>
    ),
    whiteBackground,
    centered,
  ],
}

const Template = (args) => <WithdrawETHModal {...args} />

export const Default = Template.bind({})
Default.args = {
  operatorAddress: "0xeF42ac774dD0d3519E7CBFD59F36e52038D4e255",
  availableETHInWei: "300000000000000000000",
  availableETH: "300",
  // closeModal,
  // managedGrantAddress,
  // withdrawUnbondedEth,
  // withdrawUnbondedEthAsManagedGrantee,
}