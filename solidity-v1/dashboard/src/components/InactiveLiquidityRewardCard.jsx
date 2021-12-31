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
      <LiquidityRewardCard.InactivePoolBanner link={viewPoolLink} />
      {children}
      <LiquidityRewardCard.Rewards
        isFetching={isFetching}
        lpBalance={lpBalance}
      />
      <LiquidityRewardCard.ActionButtons
        poolId={poolId}
        incentivesRemoved={true}
        wrappedTokenBalance={wrappedTokenBalance}
        lpBalance={lpBalance}
        rewardBalance={rewardBalance}
        addLpTokens={addLpTokens}
      />
    </LiquidityRewardCard>
  )
}

export default InactiveLiquidityRewardCard
