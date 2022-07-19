import React from "react"

import MetricsTile from "..//MetricsTile"
import { APY } from "../liquidity"
import { Skeleton } from "../skeletons"
import TokenAmount from "../TokenAmount"
import ResourceTooltip from "../ResourceTooltip"
import OnlyIf from "../OnlyIf"
import { TOOLTIP_DIRECTION } from "../Tooltip"

const MetricsSection = ({
  tvl,
  tvlInUSD,
  rewardRate,
  isRewardRateFetching,
  totalAllocatedRewards,
  isTotalAllocatedRewardsFetching,
  lifetimeCovered,
  isLifetimeCoveredFetching,
  classes = {},
}) => {
  return (
    <section className={`tile coverage-pool__metrics ${classes.root || ""}`}>
      <section className={`metrics__tvl ${classes.tvl || ""}`}>
        <div className="mb-1 flex row center">
          <h2 className={"h2--alt text-grey-70"}>Total Value Locked</h2>
          <ResourceTooltip
            tooltipClassName="ml-1"
            title="Total Value Locked"
            content="The total amount of KEEP deposited into the coverage pool."
            redirectLink="/coverage-pools/how-it-works"
            linkText="How it Works"
          />
        </div>
        <TokenAmount
          amount={tvl}
          amountClassName="h1 text-mint-100"
          symbolClassName="h2 text-mint-100"
          withIcon
        />
        <h3 className="tvl tvl--fiat-currency">{`$${tvlInUSD.toString()} USD`}</h3>
      </section>

      <section className={`metrics__reward-rate ${classes.rewardRate || ""}`}>
        <MetricsTile className="bg-mint-10">
          <MetricsTile.Tooltip direction={TOOLTIP_DIRECTION.TOP}>
            Estimated rate of rewards that you will receive annually.
          </MetricsTile.Tooltip>
          <APY
            apy={rewardRate}
            isFetching={isRewardRateFetching}
            className="text-mint-100"
          />
          <h5 className="text-grey-60">annual rewards rate</h5>
        </MetricsTile>
      </section>

      <section
        className={`metrics__total-rewards ${classes.totalRewards || ""}`}
      >
        <MetricsTile className="bg-mint-10">
          <MetricsTile.Tooltip direction={TOOLTIP_DIRECTION.TOP}>
            Rewards distributed from the rewards pool contract since the start
            of the pool.
          </MetricsTile.Tooltip>
          <OnlyIf condition={isTotalAllocatedRewardsFetching}>
            <Skeleton tag="h2" shining color="grey-10" />
          </OnlyIf>
          <OnlyIf condition={!isTotalAllocatedRewardsFetching}>
            <TokenAmount
              amount={totalAllocatedRewards}
              withIcon
              withSymbol={false}
              withMetricSuffix
            />{" "}
          </OnlyIf>
          <h5 className="text-grey-60">lifetime rewards</h5>
        </MetricsTile>
      </section>

      <section
        className={`metrics__lifetime-covered ${classes.lifetimeCovered || ""}`}
      >
        <MetricsTile className="bg-mint-10">
          <MetricsTile.Tooltip direction={TOOLTIP_DIRECTION.TOP}>
            Amount of KEEP used from the coverage pool to cover a loss since the
            start of the pool.
          </MetricsTile.Tooltip>
          <OnlyIf condition={isLifetimeCoveredFetching}>
            <Skeleton tag="h2" shining color="grey-10" />
          </OnlyIf>
          <OnlyIf condition={!isLifetimeCoveredFetching}>
            <TokenAmount
              amount={lifetimeCovered}
              withIcon
              withSymbol={false}
              withMetricSuffix
            />
          </OnlyIf>
          <h5 className="text-grey-60">lifetime covered</h5>
        </MetricsTile>
      </section>
    </section>
  )
}

export default MetricsSection
