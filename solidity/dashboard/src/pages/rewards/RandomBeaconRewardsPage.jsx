import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import TokenAmount from "../../components/TokenAmount"
import { SubmitButton } from "../../components/Button"
import { BeaconRewardsDetails } from "../../components/RewardsDetails"
// import StakeDropChart from "../../components/StakeDropChart"
import { useWeb3Context } from "../../components/WithWeb3Context"
import { TokenAmountSkeleton } from "../../components/skeletons"
import EmptyStatePage from "./EmptyStatePage"
import { gt } from "../../utils/arithmetics.utils"

const RandomBeaconRewardsPage = () => {
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
    <>
      <section className="rewards-overview--random-beacon">
        <section>
          <Balance
            title="Random Beacon Rewards"
            rewardsBalance={becaonRewardsBalance}
            isBalanceFetching={beaconRewardsFetching}
          />
        </section>
        <section>
          <BeaconRewardsDetails />
        </section>
        {/* For now, we decided to drop out the `StakeDropChart`
            to keep consistency.
        */}
        {/* <section className="tile">
          <StakeDropChart />
        </section> */}
      </section>
    </>
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
      {onWithdrawAll && (
        <SubmitButton
          onSubmitAction={onWithdrawAll}
          className="btn btn-primary btn-lg mt-2"
          disabled={!gt(rewardsBalance, 0)}
        >
          withdraw all
        </SubmitButton>
      )}
    </>
  )
}

RandomBeaconRewardsPage.route = {
  title: "Random Beacon",
  path: "/rewards/random-beacon",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default RandomBeaconRewardsPage
