import React from "react"
import LiquidityRewardCard from "./LiquidityRewardCard"

const ActiveLiquidityRewardCard = ({
  poolId,
  title,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
  apy,
  // Percentage of the deposited liquidity tokens in the `LPRewards` pool.
  percentageOfTotalPool,
  // Current reward balance earned in `LPRewards` contract.
  rewardBalance = "0",
  // Balance of the wrapped token.
  wrappedTokenBalance = "0",
  // Balance of wrapped token deposited in the `LPRewards` contract.
  lpBalance = "0",
  isFetching,
  wrapperClassName = "",
  addLpTokens,
  withdrawLiquidityRewards,
  isAPYFetching,
  pool,
  children,
}) => {
  return (
    <LiquidityRewardCard
      title={title}
      MainIcon={MainIcon}
      SecondaryIcon={SecondaryIcon}
      wrapperClassName={wrapperClassName}
      isFetching={isFetching}
      isAPYFetching={isAPYFetching}
      apy={apy}
      percentageOfTotalPool={percentageOfTotalPool}
      lpBalance={lpBalance}
      wrappedTokenBalance={wrappedTokenBalance}
      rewardBalance={rewardBalance}
    >
      <LiquidityRewardCard.Subtitle pool={pool} viewPoolLink={viewPoolLink} />
      <LiquidityRewardCard.Metrics />
      <LiquidityRewardCard.ActivePoolBanner
        pool={pool}
        viewPoolLink={viewPoolLink}
      />
      {children}
      <LiquidityRewardCard.Rewards />
      <LiquidityRewardCard.ActionButtons
        poolId={poolId}
        incentivesRemoved={false}
        addLpTokens={addLpTokens}
        withdrawLiquidityRewards={withdrawLiquidityRewards}
      />
    </LiquidityRewardCard>
  )
}

export default ActiveLiquidityRewardCard
