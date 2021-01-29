import React, { useMemo } from "react"
import CountUp from "react-countup"
import BigNumber from "bignumber.js"
import { Skeleton } from "../skeletons"
import {
  displayPercentageValue,
  formatPercentage,
} from "../../utils/general.utils"

const ShareOfPool = ({
  percentageOfTotalPool,
  isFetching = false,
  skeletonProps = { tag: "h2", shining: true, color: "grey-10" },
  className = "",
}) => {
  const formattedPercentageOfTotalPool = useMemo(() => {
    const bn = new BigNumber(percentageOfTotalPool)
    return bn.isLessThan(0.01) && bn.isGreaterThan(0)
      ? 0.01
      : formatPercentage(bn)
  }, [percentageOfTotalPool])

  return isFetching ? (
    <Skeleton {...skeletonProps} />
  ) : (
    <h2 className={className}>
      <CountUp
        end={formattedPercentageOfTotalPool}
        // Save previously ended number to start every new animation from it.
        preserveValue
        decimals={2}
        duration={1}
        formattingFn={displayPercentageValue}
      />
    </h2>
  )
}

export default ShareOfPool
