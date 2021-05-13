import React from "react"
import { TopUpsDataTable } from "../components/TopUpsDataTable"

export default {
  title: "TopUpsDataTable",
  component: TopUpsDataTable,
}

const mockedTopUps = [
  {
    availableTopUpAmount: "2000000000000000000000",
    createdAt: 1620803413,
    IsInUndelegation: false,
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
  },
  {
    availableTopUpAmount: "2000000000000000000000",
    createdAt: 1620803413,
    IsInUndelegation: false,
    operatorAddress: "0xc360C120Aa05bAffeE3b427cCFc7F19FBBcD9953",
  },
]

const Template = (args) => <TopUpsDataTable {...args} />

export const Default = Template.bind({})
Default.args = {
  topUps: mockedTopUps,
  commitTopUp: null,
}

