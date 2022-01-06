import React from "react"
import LiquidityRewardCard from "./LiquidityRewardCard"

const InactiveLiquidityRewardCard = ({
  poolId,
  title,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
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
      lpBalance={lpBalance}
      wrappedTokenBalance={wrappedTokenBalance}
      rewardBalance={rewardBalance}
    >
      <LiquidityRewardCard.Subtitle pool={pool} viewPoolLink={viewPoolLink} />
      <LiquidityRewardCard.InactivePoolBanner link={viewPoolLink} />
      {children}
      <LiquidityRewardCard.Rewards />
      <LiquidityRewardCard.ActionButtons
        poolId={poolId}
        incentivesRemoved={true}
        addLpTokens={addLpTokens}
        withdrawLiquidityRewards={withdrawLiquidityRewards}
      />
    </LiquidityRewardCard>
  )
}

export default InactiveLiquidityRewardCard
