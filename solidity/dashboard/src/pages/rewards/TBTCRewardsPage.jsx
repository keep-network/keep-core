import React, { useEffect, useCallback } from "react"
import { useDispatch, useSelector } from "react-redux"
// import { Link } from "react-router-dom"
import TokenAmount from "../../components/TokenAmount"
import * as Icons from "../../components/Icons"
import Tooltip from "../../components/Tooltip"
import { SubmitButton } from "../../components/Button"
// import ProgressBar from "../../components/ProgressBar"
import Timer from "../../components/Timer"
import TBTCRewardsDataTable from "../../components/TBTCRewardsDataTable"
import { ECDSARewardsHelper } from "../../utils/rewardsHelper"
import { useWeb3Address } from "../../components/WithWeb3Context"
import {
  TokenAmountSkeleton,
  DataTableSkeleton,
} from "../../components/skeletons"
import { gt } from "../../utils/arithmetics.utils"
import { LoadingOverlay } from "../../components/Loadable"
import EmptyStatePage from "./EmptyStatePage"
// import { colors } from "../../constants/colors"

const TBTCRewardsPage = () => {
  const dispatch = useDispatch()
  const yourAddress = useWeb3Address()
  const currentIntervalEndOf = ECDSARewardsHelper.intervalEndOf(
    ECDSARewardsHelper.currentInterval
  ).unix()

  const {
    ecdsaAvailableRewardsFetching,
    ecdsaAvailableRewardsBalance,
    ecdsaAvailableRewards,
    ecdsaRewardsHistory,
  } = useSelector((state) => state.rewards)

  useEffect(() => {
    dispatch({
      type: "rewards/ecdsa_fetch_rewards_data_request",
      payload: { address: yourAddress },
    })
  }, [dispatch, yourAddress])

  const withdrawRewards = useCallback(
    (submitButtonPromise) => {
      dispatch({
        type: "rewards/ecdsa_withdraw",
        payload: ecdsaAvailableRewards,
        meta: submitButtonPromise,
      })
    },
    [dispatch, ecdsaAvailableRewards]
  )

  return (
    <>
      <RewardsOverview
        balance={ecdsaAvailableRewardsBalance}
        isBalanceFetching={ecdsaAvailableRewardsFetching}
        withdrawRewards={withdrawRewards}
      />
      <section className="tile rewards-countdown">
        <h2 className="h2--alt">
          Next rewards release:&nbsp;
          <Timer targetInUnix={currentIntervalEndOf} />
        </h2>
      </section>
      <LoadingOverlay
        isFetching={ecdsaAvailableRewardsFetching}
        skeletonComponent={<DataTableSkeleton columns={4} />}
      >
        <section className="tile rewards-history">
          <TBTCRewardsDataTable data={ecdsaRewardsHistory} />
        </section>
      </LoadingOverlay>
    </>
  )
}

const RewardsOverview = ({ balance, isBalanceFetching, withdrawRewards }) => {
  const remainingPeriods =
    ECDSARewardsHelper.intervals - ECDSARewardsHelper.currentInterval

  return (
    <section className="tile rewards__overview--tbtc">
      <div className="rewards__overview__balance">
        <h2 className="h2--alt text-grey-70 mb-1">tBTC Rewards</h2>
        {isBalanceFetching ? (
          <TokenAmountSkeleton />
        ) : (
          <TokenAmount amount={balance} withIcon withMetricSuffix />
        )}
      </div>
      <div className="rewards__overview__period">
        <h5 className="text-grey-70">current rewards period</h5>
        <span className="rewards-period__date">
          {ECDSARewardsHelper.periodOf(ECDSARewardsHelper.currentInterval)}
        </span>
        <div className="rewards-period__remaining-periods">
          <Icons.Time width="16" height="16" className="time-icon--grey-30" />
          <span>{remainingPeriods} rewards periods remaining&nbsp;</span>
          <Tooltip simple delay={0} triggerComponent={Icons.MoreInfo}>
            Rewards are distributed for a limited time.
          </Tooltip>
        </div>
      </div>
      <div className="rewards__overview__withdraw">
        <SubmitButton
          className="btn btn-primary btn-lg w-100"
          onSubmitAction={withdrawRewards}
          disabled={!gt(balance, 0)}
        >
          withdraw all
        </SubmitButton>
        <div className="text-validation mt-1 text-center">
          Beneficiary account receives rewards.
        </div>
      </div>
    </section>
  )
}

// const RewardsBoost = () => {
//   return (
//     <section className="tile rewards-calc__boost">
//       <div>
//         <div className="boost__value">
//           <Icons.Plus width={48} height={48} />
//           <div className="ml-1">
//             <h3 className="text-grey-70">2.5x Boost</h3>
//             <h4 className="text-mint-100">500k KEEP</h4>
//           </div>
//         </div>
//         <ProgressBar value={10} total={100} color={colors.primary}>
//           <ProgressBar.Inline className="mb-0" />
//         </ProgressBar>
//         <div className="text-grey-40">500k KEEP</div>
//         <Link to="/delegations" className="btn btn-secondary mt-1 w-100">
//           stake keep
//         </Link>
//       </div>
//       <div>
//         <p>
//           Up to 3M KEEP staked per operator counts toward your reward boost.
//           Your boost multiplies the rewards you receive from bonding ETH.
//         </p>
//         <p>Ahigher boost means more rewards.</p>
//       </div>
//     </section>
//   )
// }

// const LockedETH = () => {
//   return (
//     <section className="tile rewards-calc__eth-locked">
//       <div className="eth-locked__value">
//         <Icons.ETH className="" />
//         <div className="ml-1">
//           <h3 className="text-grey-70">ETH Locked</h3>
//           <h4 className="text-mint-100">1000 ETH</h4>
//         </div>
//       </div>
//       <ProgressBar value={10} total={100} color={colors.primary}>
//         <ProgressBar.Inline className="mb-0" />
//       </ProgressBar>
//       <div className="text-grey-40">1000 ETH Locked</div>
//       <Link to="/applications/tbtc" className="btn btn-secondary mt-1 w-100">
//         lock eth
//       </Link>
//     </section>
//   )
// }

// const APY = () => {
//   return (
//     <section className="tile rewards-calc__apy">
//       <h2 className="apy__title">Your APY</h2>
//       <h1 className="apy__value">252%</h1>
//     </section>
//   )
// }

TBTCRewardsPage.route = {
  title: "tBTC",
  path: "/rewards/tbtc",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default TBTCRewardsPage
