import React from "react"
import { Skeleton } from "../skeletons"
import CountUp from "react-countup"

export const RewardMultiplier = ({
  rewardMultiplier,
  isFetching = false,
  skeletonProps = { tag: "h2", shining: true, color: "grey-10" },
  className = "",
}) => {
  return isFetching ? (
    <Skeleton {...skeletonProps} />
  ) : (
    <h2 className={` ${className} liquidity__info-tile__title text-mint-100`}>
      <CountUp
        end={rewardMultiplier}
        preserveValue
        decimals={1}
        duration={1}
        suffix={"x"}
      />
    </h2>
  )
}

RewardMultiplier.TooltipContent = () => {
  // TODO: content of the tooltip
  return <>Reward multiplier tooltip</>
}

export default RewardMultiplier
