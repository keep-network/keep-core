import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import TokenAmount from "../../components/TokenAmount"
// import { SubmitButton } from "../../components/Button"
import BeaconRewardsDetails from "../../components/BeaconRewardsDetails"
import StakeDropChart from "../../components/StakeDropChart"
import { useWeb3Context } from "../../components/WithWeb3Context"
import { TokenAmountSkeleton } from "../../components/skeletons"

const RewardsOverviewPage = () => {
  const dispatch = useDispatch()
  const { yourAddress } = useWeb3Context()
  const { beaconRewardsFetching, becaonRewardsBalance } = useSelector(
    (state) => state.rewards
  )

  useEffect(() => {
    dispatch({
      type: "rewards/beacon_fetch_distributed_rewards_request",
      payload: yourAddress,
    })
  }, [dispatch, yourAddress])

  return (
    <section className="rewards-overview--random-beacon">
      <section className="tile">
        <Balance
          title="Random Beacon Rewards"
          rewardsBalance={becaonRewardsBalance}
          isBalanceFetching={beaconRewardsFetching}
        />
        <section className="mt-2">
          <BeaconRewardsDetails />
        </section>
      </section>
      <section className="tile">
        <StakeDropChart />
      </section>
    </section>
  )
}

const Balance = ({
  title,
  rewardsBalance,
  isBalanceFetching,
  onWithdrawAll,
}) => {
  return (
    <>
      <h2 className="h2--alt mb-1">{title}</h2>
      {isBalanceFetching ? (
        <TokenAmountSkeleton />
      ) : (
        <TokenAmount amount={rewardsBalance} currencySymbol="KEEP" />
      )}

      {/* <div className="flex column wrap ">
        <SubmitButton onClick={onWithdrawAll}>withdraw all</SubmitButton>
        <span className="text-validation">
          The beneficiary account receives all withdrawn rewards.
        </span>
      </div> */}
    </>
  )
}

const EmptyStatePage = () => <>connect wallet</>

RewardsOverviewPage.route = {
  title: "Overview",
  path: "/rewards/overview",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default RewardsOverviewPage
