import React from "react"
import MetricsTile from "../components/MetricsTile"
import RewardMultiplier from "../components/liquidity/RewardMultiplier"
import { APY, ShareOfPool } from "../components/liquidity"
import { storiesOf } from "@storybook/react"
import centered from "@storybook/addon-centered/react"

storiesOf("MetricsTile", module).addDecorator(centered)

export default {
  title: "MetricsTile",
  component: MetricsTile,
  subcomponents: { RewardMultiplier },
  decorators: [
    (Story) => (
      <section
        className={`keep-only-pool__overview__info-tiles liquidity__info`}
      >
        <Story />
      </section>
    ),
  ],
}

const Template = (args) => <MetricsTile {...args} />

export const ForRewardMultiplier = Template.bind({})
ForRewardMultiplier.args = {
  children: (
    <>
      <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
        <RewardMultiplier.TooltipContent />
      </MetricsTile.Tooltip>
      <RewardMultiplier
        rewardMultiplier={2.6}
        className="liquidity__info-tile__title text-mint-100"
      />
      <h6>reward multiplier</h6>
    </>
  ),
  className: "liquidity__info-tile bg-mint-10",
}

export const ForAPY = Template.bind({})
ForAPY.args = {
  children: (
    <>
      <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
        <APY.TooltipContent />
      </MetricsTile.Tooltip>
      <APY apy={0.999} className="liquidity__info-tile__title text-mint-100" />
      <h6>Estimate of pool apy</h6>
    </>
  ),
  className: "liquidity__info-tile bg-mint-10",
}

export const ForShareOfPool = Template.bind({})
ForShareOfPool.args = {
  children: (
    <>
      <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
        <ShareOfPool.TooltipContent />
      </MetricsTile.Tooltip>
      <ShareOfPool
        percentageOfTotalPool={27}
        className="liquidity__info-tile__title text-mint-100"
      />
      <h6>Your share of POOL</h6>
    </>
  ),
  className: "liquidity__info-tile bg-mint-10",
}

export const GrayedOut = Template.bind({})
GrayedOut.decorators = [
  (Story) => (
    <section
      className={`keep-only-pool__overview__info-tiles liquidity__info--locked`}
    >
      <Story />
    </section>
  ),
]
GrayedOut.args = {
  children: (
    <>
      <MetricsTile.Tooltip className="liquidity__info-tile__tooltip">
        <RewardMultiplier.TooltipContent />
      </MetricsTile.Tooltip>
      <RewardMultiplier
        rewardMultiplier={0}
        className="liquidity__info-tile__title text-mint-100"
      />
      <h6>reward multiplier</h6>
    </>
  ),
  className: "liquidity__info-tile bg-mint-10",
}
