import React, { useMemo } from "react"
import CountUp from "react-countup"
import BigNumber from "bignumber.js"
import { Skeleton } from "../skeletons"
import {
  displayPercentageValue,
  formatPercentage,
} from "../../utils/general.utils"

export const APY = ({
  apy,
  isFetching = false,
  skeletonProps = { tag: "h2", shining: true, color: "grey-10" },
  className = "",
}) => {
  const formattedApy = useMemo(() => {
    const bn = new BigNumber(apy).multipliedBy(100)
    if (bn.isEqualTo(Infinity)) {
      return Infinity
    } else if (bn.isLessThan(0.01) && bn.isGreaterThan(0)) {
      return 0.01
    } else if (bn.isGreaterThan(999)) {
      return 999
    }

    return formatPercentage(bn)
  }, [apy])

  return isFetching ? (
    <Skeleton {...skeletonProps} />
  ) : (
    <h2 className={className}>
      {formattedApy === Infinity ? (
        <span>&#8734;</span>
      ) : (
        <CountUp
          end={formattedApy}
          // Save previously ended number to start every new animation from it.
          preserveValue
          decimals={2}
          duration={1}
          formattingFn={displayPercentageValue}
        />
      )}
    </h2>
  )
}

APY.TooltipContent = () => {
  return (
    <>
      Pool APY is calculated using the&nbsp;
      <a
        target="_blank"
        rel="noopener noreferrer"
        href={"https://thegraph.com/explorer/subgraph/uniswap/uniswap-v2"}
        className="text-white text-link"
      >
        Uniswap subgraph API
      </a>
      &nbsp;to fetch the total pool value and KEEP token in USD.
    </>
  )
}

export default APY
