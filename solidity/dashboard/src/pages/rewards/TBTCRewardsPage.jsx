import React from "react"
import TokenAmount from "../../components/TokenAmount"
import * as Icons from "../../components/Icons"
import Tooltip from "../../components/Tooltip"
import { Link } from "react-router-dom"
import { SubmitButton } from "../../components/Button"
import ProgressBar from "../../components/ProgressBar"
import Timer from "../../components/Timer"
import TBTCRewardsDataTable from "../../components/TBTCRewardsDataTable"
import { colors } from "../../constants/colors"
import { ECDSARewardsHelper } from "../../utils/rewardsHelper";

const TBTCRewardsPage = () => {
  const currentIntervalEndOf = ECDSARewardsHelper.intervalEndOf(
    ECDSARewardsHelper.currentInterval
  ).unix()

  return (
    <>
      <RewardsOverview />
      <section className="tile rewards-countdown">
        <h2 className="h2--alt">
          Next rewards release:&nbsp;
          <Timer targetInUnix={currentIntervalEndOf} />
        </h2>
      </section>
      <section className="rewards-calc">
        <RewardsBoost />
        <LockedETH />
        <APY />
      </section>
      <section className="rewards-history">
        <TBTCRewardsDataTable />
      </section>
    </>
  )
}

const RewardsOverview = () => {
  return (
    <section className="tile rewards__overview--tbtc">
      <div className="rewards__overview__balance">
        <h2 className="h2--alt text-grey-70 mb-1">tBTC Rewards</h2>
        <TokenAmount amount="0" currencySymbol="KEEP" />
      </div>
      <div className="rewards__overview__period">
        <h5 className="text-grey-70">current rewards period</h5>
        <span className="rewards-period__date">11/15/2020 - 12/15/2020</span>
        <div className="rewards-period__remaining-periods">
          <Icons.Time width="16" height="16" className="time-icon--grey-30" />
          <span>84 rewards periods remaining&nbsp;</span>
          <Tooltip simple delay={0} triggerComponent={Icons.MoreInfo}>
            content
          </Tooltip>
        </div>
      </div>
      <div className="rewards__overview__withdraw">
        <SubmitButton
          className="btn btn-primary btn-lg"
          onSubmitAction={() => console.log("submit btn")}
        >
          withdraw all
        </SubmitButton>
        <div className="text-validation mt-1">
          Beneficiary account receives rewards.
        </div>
      </div>
    </section>
  )
}

const RewardsBoost = () => {
  return (
    <section className="tile rewards-calc__boost">
      <div>
        <div className="boost__value">
          <Icons.Plus width={48} height={48} />
          <div className="ml-1">
            <h3 className="text-grey-70">2.5x Boost</h3>
            <h4 className="text-mint-100">500k KEEP</h4>
          </div>
        </div>
        <ProgressBar value={10} total={100} color={colors.primary}>
          <ProgressBar.Inline className="mb-0" />
        </ProgressBar>
        <div className="text-grey-40">500k KEEP</div>
        <Link to="/delegations" className="btn btn-secondary mt-1 w-100">
          stake keep
        </Link>
      </div>
      <div>
        <p>
          Up to 3M KEEP staked per operator counts toward your reward boost.
          Your boost multiplies the rewards you receive from bonding ETH.
        </p>
        <p>Ahigher boost means more rewards.</p>
      </div>
    </section>
  )
}

const LockedETH = () => {
  return (
    <section className="tile rewards-calc__eth-locked">
      <div className="eth-locked__value">
        <Icons.ETH className="" />
        <div className="ml-1">
          <h3 className="text-grey-70">ETH Locked</h3>
          <h4 className="text-mint-100">1000 ETH</h4>
        </div>
      </div>
      <ProgressBar value={10} total={100} color={colors.primary}>
        <ProgressBar.Inline className="mb-0" />
      </ProgressBar>
      <div className="text-grey-40">1000 ETH Locked</div>
      <Link to="/applications/tbtc" className="btn btn-secondary mt-1 w-100">
        lock eth
      </Link>
    </section>
  )
}

const APY = () => {
  return (
    <section className="tile rewards-calc__apy">
      <h2 className="apy__title">Your APY</h2>
      <h1 className="apy__value">252%</h1>
    </section>
  )
}

TBTCRewardsPage.route = {
  title: "tBTC",
  path: "/rewards/tbtc",
  exact: true,
  withConnectWalletGuard: false,
  // TODO: empty state page component
  //   emptyStateComponent: Component,
}

export default TBTCRewardsPage
