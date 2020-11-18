import React, { useEffect } from "react"
import { useDispatch } from "react-redux"
import TokenAmount from "../../components/TokenAmount"
// import { SubmitButton } from "../../components/Button"
import BeaconRewardsDetails from "../../components/BeaconRewardsDetails"
import StakeDropChart from "../../components/StakeDropChart"
import { useWeb3Context } from "../../components/WithWeb3Context"

const RewardsOverviewPage = () => {
  const dispatch = useDispatch()
  const { yourAddress } = useWeb3Context()

  useEffect(() => {
    dispatch({
      type: "rewards/beacon_fetch_distributed_rewards_request",
      payload: yourAddress,
    })
  }, [dispatch, yourAddress])

  return (
    <section className="rewards-overview--random-beacon">
      <section className="tile">
        <Balance title="Random Beacon Rewards" rewardsBalance={0} />
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

const Balance = ({ title, rewardsBalance, onWithdrawAll }) => {
  return (
    <>
      <h2 className="h2--alt mb-1">{title}</h2>
      <TokenAmount amount={rewardsBalance} currencySymbol="KEEP" />
      {/* <div className="flex column wrap ">
        <SubmitButton onClick={onWithdrawAll}>withdraw all</SubmitButton>
        <span className="text-validation">
          The beneficiary account receives all withdrawn rewards.
        </span>
      </div> */}
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
