import React from "react"

import MetricsTile from "../../components/MetricsTile"
import { APY } from "../../components/liquidity"
import { Skeleton } from "../../components/skeletons"
import TokenAmount from "../../components/TokenAmount"
import ResourceTooltip from "../ResourceTooltip"

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
        <h2 className="h2--alt text-grey-70 mb-1">
          Total Value Locked
          <ResourceTooltip
            tooltipClassName="ml-1"
            title="Total Value Locked"
            content="The total amount of KEEP deposited into the coverage pool."
            redirectLink="/coverage-pools/how-it-works"
            linkText="How it Works"
          />
        </h2>
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
          <MetricsTile.Tooltip direction="top">
            The rate of rewards that you will receive annually.
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
          <MetricsTile.Tooltip direction="top">
            Rewards distributed from the rewards pool contract since the start
            of the pool.
          </MetricsTile.Tooltip>
          {isTotalAllocatedRewardsFetching ? (
            <Skeleton tag="h2" shining color="grey-10" />
          ) : (
            <TokenAmount
              amount={totalAllocatedRewards}
              withIcon
              withSymbol={false}
              withMetricSuffix
            />
          )}
          <h5 className="text-grey-60">lifetime rewards</h5>
        </MetricsTile>
      </section>

      <section
        className={`metrics__lifetime-covered ${classes.lifetimeCovered || ""}`}
      >
        <MetricsTile className="bg-mint-10">
          <MetricsTile.Tooltip direction="top">
            Amount of KEEP used from the coverage pool to cover a loss since the
            start of the pool.
          </MetricsTile.Tooltip>
          {isLifetimeCoveredFetching ? (
            <Skeleton tag="h2" shining color="grey-10" />
          ) : (
            <TokenAmount
              amount={lifetimeCovered}
              withIcon
              withSymbol={false}
              withMetricSuffix
            />
          )}
          <h5 className="text-grey-60">lifetime covered</h5>
        </MetricsTile>
      </section>
    </section>
  )
}

export default MetricsSection
