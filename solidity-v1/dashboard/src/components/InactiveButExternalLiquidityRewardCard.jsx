import React from "react"
import LiquidityRewardCard from "./LiquidityRewardCard"
import * as Icons from "./Icons"

const InactiveButExternalLiquidityRewardCard = ({
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
  inactivePoolBannerProps = {
    icon: Icons.Rewards,
    title: "Incentives moved to T",
    description:
      "You can still withdraw rewards that you already earned. Click `Go to pool` and deposit your tokens in the Saddle dApp to earn incentives in multiple tokens.",
    link: "https://forum.keep.network/t/repurpose-saddle-tbtc-pool-liquidity-incentives-and-move-incentives-to-t/404",
    linkText: "More info",
  },
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
      <LiquidityRewardCard.ActivePoolBanner
        pool={pool}
        viewPoolLink={viewPoolLink}
        userInfoBannerProps={inactivePoolBannerProps}
      />
      {children}
      <LiquidityRewardCard.Rewards />
      <LiquidityRewardCard.GoToPoolButton viewPoolLink={viewPoolLink} />
      <LiquidityRewardCard.ActionButtons
        poolId={poolId}
        incentivesRemoved={true}
        addLpTokens={addLpTokens}
        withdrawLiquidityRewards={withdrawLiquidityRewards}
      />
    </LiquidityRewardCard>
  )
}

export default InactiveButExternalLiquidityRewardCard
