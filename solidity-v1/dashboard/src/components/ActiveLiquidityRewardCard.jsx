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
    >
      <LiquidityRewardCard.Subtitle pool={pool} viewPoolLink={viewPoolLink} />
      <LiquidityRewardCard.Metrics
        apy={apy}
        isFetching={isFetching}
        isAPYFetching={isAPYFetching}
        lpBalance={lpBalance}
        percentageOfTotalPool={percentageOfTotalPool}
      />
      <LiquidityRewardCard.ActivePoolBanner
        pool={pool}
        viewPoolLink={viewPoolLink}
        lpBalance={lpBalance}
      />
      {children}
      <LiquidityRewardCard.Rewards
        isFetching={isFetching}
        lpBalance={lpBalance}
      />
      <LiquidityRewardCard.ActionButtons
        poolId={poolId}
        incentivesRemoved={false}
        wrappedTokenBalance={wrappedTokenBalance}
        lpBalance={lpBalance}
        rewardBalance={rewardBalance}
        addLpTokens={addLpTokens}
        withdrawLiquidityRewards={withdrawLiquidityRewards}
      />
    </LiquidityRewardCard>
  )
}

export default ActiveLiquidityRewardCard
