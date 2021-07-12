import React from "react"
import centered from "@storybook/addon-centered/react"
import RewardsStatus from "../components/RewardsStatus"
import { PENDING_STATUS, REWARD_STATUS } from "../constants/constants"

const mockedTransactionHash =
  "6146ccf6a66d994f7c363db875e31ca35581450a4bf6d3be6cc9ac79233a69d0"

export default {
  title: "RewardsStatus",
  component: RewardsStatus,
  decorators: [centered],
}

const Template = (args) => <RewardsStatus {...args} />

export const Pending = Template.bind({})
Pending.args = {
  status: PENDING_STATUS,
  transactionHash: mockedTransactionHash,
}

export const Available = Template.bind({})
Available.args = {
  status: REWARD_STATUS.AVAILABLE,
  transactionHash: mockedTransactionHash,
}

export const Accumulating = Template.bind({})
Accumulating.args = {
  status: REWARD_STATUS.ACCUMULATING,
  transactionHash: mockedTransactionHash,
}

export const Withdrawn = Template.bind({})
Withdrawn.args = {
  status: REWARD_STATUS.WITHDRAWN,
  transactionHash: mockedTransactionHash,
}
