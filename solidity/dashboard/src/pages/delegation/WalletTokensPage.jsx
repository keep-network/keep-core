import React from "react"
import { useSelector, useDispatch } from "react-redux"
import { LoadingOverlay } from "../../components/Loadable"
import DelegatedTokensTable from "../../components/DelegatedTokensTable"
import Undelegations from "../../components/Undelegations"
import { DataTableSkeleton } from "../../components/skeletons"
import EmptyStateComponent from "./EmptyStatePage"
import TokenAmount from "../../components/TokenAmount"
import { colors } from "../../constants/colors"
import DelegateStakeForm from "../../components/DelegateStakeForm"
import { useModal } from "../../hooks/useModal"
import moment from "moment"
import ProgressBar from "../../components/ProgressBar"

const confirmationModalOptions = (initializationPeriod) => ({
  modalOptions: { title: "Initiate Delegation" },
  title: "You’re about to delegate stake.",
  subtitle: `You’re delegating KEEP tokens. You will be able to cancel the delegation for up to ${moment()
    .add(initializationPeriod, "seconds")
    .fromNow(true)}. After that time, you can undelegate your stake.`,
  btnText: "delegate",
  confirmationText: "DELEGATE",
})

const WalletTokensPage = () => {
  const dispatch = useDispatch()
  const { openConfirmationModal } = useModal()

  const {
    minimumStake,
    initializationPeriod,
    undelegationPeriod,
    delegations,
    undelegations,
    isDelegationDataFetching,
  } = useSelector((state) => state.staking)

  const keepToken = useSelector((state) => state.keepTokenBalance)

  const handleSubmit = async (values, meta) => {
    await openConfirmationModal(confirmationModalOptions(initializationPeriod))
    dispatch({
      type: "staking/delegate_request",
      payload: {
        ...values,
        amount: values.stakeTokens,
      },
      meta,
    })
  }

  return (
    <>
      <section className="wallet-page__overview-layout">
        <section className="tile wallet-page__overview__balance">
          <h4 className="mb-1">Wallet Balance</h4>
          <TokenAmount
            amount="100000000000000000000000"
            currencySymbol="KEEP"
          />
        </section>
        <section className="tile wallet-page__overview__staked-tokens">
          <h4 className="mb-2">Tokens Staked</h4>
          <ProgressBar
            value={20}
            total={100}
            color={colors.mint80}
            bgColor={colors.mint20}
          >
            <div className="circular-progress-bar-percentage-label-wrapper">
              <ProgressBar.Circular radius={82} barWidth={16} />
              <ProgressBar.PercentageLabel text="Staked" />
            </div>
            <ProgressBar.Legend leftValueLabel="Unstaked" valueLabel="Staked" />
          </ProgressBar>
        </section>
        <section className="tile wallet-page__overview__stake-form">
          {/* TODO add tooltip. PR is in progress https://github.com/keep-network/keep-core/pull/2135  */}
          <h3 className="mb-1">Stake Wallet Tokens</h3>
          <DelegateStakeForm
            onSubmit={handleSubmit}
            minStake={minimumStake}
            availableToStake={keepToken.value}
          />
        </section>
      </section>
      <section>
        <h2 className="h2--alt text-grey-60 mb-2">Activity</h2>
        <LoadingOverlay
          isFetching={isDelegationDataFetching}
          skeletonComponent={<DataTableSkeleton />}
        >
          <DelegatedTokensTable
            delegatedTokens={delegations}
            //   cancelStakeSuccessCallback={cancelStakeSuccessCallback}
            keepTokenBalance={keepToken.value}
            undelegationPeriod={undelegationPeriod}
          />
        </LoadingOverlay>
        <LoadingOverlay
          isFetching={isDelegationDataFetching}
          skeletonComponent={<DataTableSkeleton />}
        >
          <Undelegations undelegations={undelegations} />
        </LoadingOverlay>
      </section>
    </>
  )
}

WalletTokensPage.route = {
  title: "Wallet Tokens",
  path: "/delegation/wallet",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export { WalletTokensPage }
