import React from "react"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"
import { BondingSection } from "../components/BondingSection"

storiesOf("BondingSection", module).addDecorator(centered)

const mockData = [
  {
    operatorAddress: "0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756",
    isWithdrawableForOperator: true,
    stakeAmount: "20000000000000000000000000",
    bondedETH: "0",
    bondedETHInWei: "0",
    availableETH: "0",
    availableETHInWei: "0",
  },
  {
    operatorAddress: "0xf119557AC33585405467135eC9A343DCDb047517",
    isWithdrawableForOperator: true,
    stakeAmount: "20000000000000000000000000",
    bondedETH: "0",
    bondedETHInWei: "0",
    availableETH: "50",
    availableETHInWei: "50000000000000000000",
  },
  {
    operatorAddress: "0xd2C6168Fd106908Df71Ab639f8b7e2F971Ab8205",
    isWithdrawableForOperator: true,
    stakeAmount: "20000000000000000000000000",
    bondedETH: "0",
    bondedETHInWei: "0",
    availableETH: "50",
    availableETHInWei: "50000000000000000000",
  },
  {
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
    isWithdrawableForOperator: true,
    stakeAmount: "20000000000000000000000000",
    bondedETH: "0",
    bondedETHInWei: "0",
    availableETH: "50",
    availableETHInWei: "50000000000000000000",
  },
  {
    operatorAddress: "0xCDAfb5A23A1F1c6f80706Cc101BCcf4b9A1A3e3B",
    isWithdrawableForOperator: true,
    stakeAmount: "20000000000000000000000000",
    bondedETH: "0",
    bondedETHInWei: "0",
    availableETH: "50",
    availableETHInWei: "50000000000000000000",
  },
]

export default {
  title: "BondingSection",
  component: BondingSection,
}

const Template = (args) => <BondingSection {...args} />

export const Empty = Template.bind({})
Empty.args = { data: [] }

export const WithData = Template.bind({})
WithData.args = { data: mockData }
