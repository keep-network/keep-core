import React from "react"
import LiquidityRewardCard from "./LiquidityRewardCard"

const ExternalPoolLiquidityRewardCard = ({
  title,
  MainIcon,
  SecondaryIcon,
  viewPoolLink,
  wrapperClassName = "",
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
      <LiquidityRewardCard.ActivePoolBanner
        pool={pool}
        viewPoolLink={viewPoolLink}
        userInfoBannerProps={{
          description:
            "Deposit your TBTC into the mStable pool to earn with low impermanent loss risk.",
          linkText: "",
        }}
      />
      {children}
      <LiquidityRewardCard.GoToPoolButton viewPoolLink={viewPoolLink} />
    </LiquidityRewardCard>
  )
}

export default ExternalPoolLiquidityRewardCard
