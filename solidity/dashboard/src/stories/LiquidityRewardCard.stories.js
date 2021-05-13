import React from "react"
import LiquidityRewardCard from "../components/LiquidityRewardCard"
import { LIQUIDITY_REWARD_PAIRS } from "../constants/constants"
import * as Icons from "../components/Icons"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("LiquidityRewardCard", module).addDecorator(centered)

export default {
  title: "LiquidityRewardCard",
  component: LiquidityRewardCard,
  argTypes: {
    addLpTokens: {
      action: "addLpTokens clicked",
    },
    withdrawLiquidityRewards: {
      action: "withdrawLiquidityRewards clicked",
    },
  },
}

const Template = (args) => <LiquidityRewardCard {...args} />

export const KEEP_ETH = Template.bind({})
KEEP_ETH.args = {
  title: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label,
  liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName,
  MainIcon: Icons.KeepBlackGreen,
  SecondaryIcon: Icons.EthToken,
  viewPoolLink: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink,
  apy: 0.9,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool: 90,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance: "300000000000000000000",
  // Balance of the wrapped token.
  wrappedTokenBalance: "300000000000000000000",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance: "300000000000000000000",
  lpTokenBalance: {
    token0: "300000000000000000000",
    token1: "300000000000000000000",
  },
  lpTokens: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.lpTokens,
  isFetching: false,
  wrapperClassName: "keep-eth",
  isAPYFetching: false,
  pool: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.pool,
}

export const KEEP_TBTC = Template.bind({})
KEEP_TBTC.args = {
  title: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.label,
  liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.contractName,
  MainIcon: Icons.KeepBlackGreen,
  SecondaryIcon: Icons.TBTC,
  viewPoolLink: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.viewPoolLink,
  apy: 0.9,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool: 90,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance: "300000000000000000000",
  // Balance of the wrapped token.
  wrappedTokenBalance: "300000000000000000000",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance: "300000000000000000000",
  lpTokenBalance: {
    token0: "300000000000000000000",
    token1: "300000000000000000000",
  },
  lpTokens: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.lpTokens,
  isFetching: false,
  wrapperClassName: "keep-tbtc",
  isAPYFetching: false,
  pool: LIQUIDITY_REWARD_PAIRS.KEEP_TBTC.pool,
}

export const TBTC_ETH = Template.bind({})
TBTC_ETH.args = {
  title: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.label,
  liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.contractName,
  MainIcon: Icons.TBTC,
  SecondaryIcon: Icons.EthToken,
  viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.viewPoolLink,
  apy: 0.9,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool: 90,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance: "300000000000000000000",
  // Balance of the wrapped token.
  wrappedTokenBalance: "300000000000000000000",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance: "300000000000000000000",
  lpTokenBalance: {
    token0: "300000000000000000000",
    token1: "300000000000000000000",
  },
  lpTokens: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.lpTokens,
  isFetching: false,
  wrapperClassName: "tbtc-eth",
  isAPYFetching: false,
  pool: LIQUIDITY_REWARD_PAIRS.TBTC_ETH.pool,
}

export const TBTC_SADDLE = Template.bind({})
TBTC_SADDLE.args = {
  title: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.label,
  liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.contractName,
  MainIcon: Icons.TBTC,
  SecondaryIcon: Icons.Saddle,
  viewPoolLink: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.viewPoolLink,
  apy: 0.9,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool: 90,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance: "300000000000000000000",
  // Balance of the wrapped token.
  wrappedTokenBalance: "300000000000000000000",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance: "300000000000000000000",
  lpTokenBalance: {
    token0: "300000000000000000000",
    token1: "300000000000000000000",
  },
  lpTokens: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.lpTokens,
  isFetching: false,
  wrapperClassName: "tbtc-saddle",
  isAPYFetching: false,
  pool: LIQUIDITY_REWARD_PAIRS.TBTC_SADDLE.pool,
}

export const IsFetching = Template.bind({})
IsFetching.args = {
  title: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.label,
  liquidityPairContractName: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.contractName,
  MainIcon: Icons.KeepBlackGreen,
  SecondaryIcon: Icons.EthToken,
  viewPoolLink: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.viewPoolLink,
  apy: 0.9,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool: 90,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance: "300000000000000000000",
  // Balance of the wrapped token.
  wrappedTokenBalance: "300000000000000000000",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance: "300000000000000000000",
  lpTokenBalance: {
    token0: "300000000000000000000",
    token1: "300000000000000000000",
  },
  lpTokens: [
    {
      tokenName: "KEEP",
      iconName: "KeepBlackGreen",
    },
    {
      tokenName: "TBTC",
      iconName: "TBTC",
    },
  ],
  isFetching: true,
  wrapperClassName: "keep-eth",
  isAPYFetching: false,
  pool: LIQUIDITY_REWARD_PAIRS.KEEP_ETH.pool,
}
