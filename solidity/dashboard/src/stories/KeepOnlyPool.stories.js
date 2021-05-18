import React from "react"
import KeepOnlyPool from "../components/KeepOnlyPool"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"

export default {
  title: "KeepOnlyPool",
  component: KeepOnlyPool,
}

const Template = (args) => <KeepOnlyPool {...args} />

export const Default = Template.bind({})
Default.args = {
  apy: 0.9,
  lpBalance: "300000000000000000000",
  rewardBalance: "300000000000000000000",
  wrappedTokenBalance: "300000000000000000000",
  isFetching: false,
  isAPYFetching: false,
  addLpTokens: null,
  withdrawLiquidityRewards: null,
  liquidityContractName: LIQUIDITY_REWARD_PAIRS.KEEP_ONLY.contractName,
  pool: LIQUIDITY_REWARD_PAIRS.KEEP_ONLY.pool,
  rewardMultiplier: 1.4,
}
