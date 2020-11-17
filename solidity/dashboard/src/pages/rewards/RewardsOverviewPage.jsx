import React from "react"
import TokenAmount from "../../components/TokenAmount"
import { SubmitButton } from "../../components/Button"
import BeaconRewardsDetails from "../../components/BeaconRewardsDetails"
import StakeDropChart from "../../components/StakeDropChart"

const RewardsOverviewPage = () => {
  return (
    <section className="rewards-overview--random-beacon">
      <section className="tile">
        <Balance
          title="Random Beacon Rewards"
          rewardsBalance={0}
          onWithdrawAll={() => console.log("on withdraw btn")}
        />
      </section>
      <section className="tile">
        <BeaconRewardsDetails />
      </section>
      <section className="tile">
        <StakeDropChart />
      </section>
    </section>
  )
}

const Balance = ({ title, onWithdrawAll }) => {
  return (
    <>
      <h2 className="h2--alt">{title}</h2>
      <TokenAmount amount={0} currencySymbol="KEEP" />
      <div className="flex column wrap ">
        <SubmitButton onClick={onWithdrawAll}>withdraw all</SubmitButton>
        <span className="text-validation">
          The beneficiary account receives all withdrawn rewards.
        </span>
      </div>
    </>
  )
}

RewardsOverviewPage.route = {
  title: "Overview",
  path: "/rewards/overview",
  exact: true,
  withConnectWalletGuard: false,
}

export default RewardsOverviewPage
