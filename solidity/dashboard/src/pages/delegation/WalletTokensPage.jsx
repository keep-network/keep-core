import React, { useMemo } from "react"
import { useSelector } from "react-redux"
import EmptyStateComponent from "./EmptyStatePage"
import TokenAmount from "../../components/TokenAmount"
import { colors } from "../../constants/colors"
import DelegateStakeForm from "../../components/DelegateStakeForm"
import ProgressBar from "../../components/ProgressBar"
import { DelegationPageWrapper } from "./index"
import { add } from "../../utils/arithmetics.utils"
import { displayAmountWithMetricSuffix } from "../../utils/token.utils"
import DelegationOverview from "../../components/DelegationOverview"
import ResourceTooltip from "../../components/ResourceTooltip"
import resourceTooltipProps from "../../constants/tooltips"

const filterByOwned = (delegation) => !delegation.grantId

const WalletTokensPageComponent = ({ onSubmitDelegateStakeForm }) => {
  const {
    minimumStake,
    initializationPeriod,
    undelegationPeriod,
    delegations,
    undelegations,
    isDelegationDataFetching,
    ownedTokensDelegationsBalance,
    ownedTokensUndelegationsBalance,
    areTopUpsFetching,
    topUps,
  } = useSelector((state) => state.staking)

  const keepToken = useSelector((state) => state.keepTokenBalance)

  const ownedDelegations = useMemo(() => {
    return delegations.filter(filterByOwned)
  }, [delegations])

  const ownedUndelegations = useMemo(() => {
    return undelegations.filter(filterByOwned)
  }, [undelegations])

  const totalOwnedStakedBalance = useMemo(() => {
    return add(
      ownedTokensDelegationsBalance,
      ownedTokensUndelegationsBalance
    ).toString()
  }, [ownedTokensDelegationsBalance, ownedTokensUndelegationsBalance])

  const totalBalance = useMemo(() => {
    return add(totalOwnedStakedBalance, keepToken.value).toString()
  }, [keepToken.value, totalOwnedStakedBalance])

  return (
    <>
      <section className="wallet-page__overview-layout">
        <section className="tile wallet-page__overview__balance">
          <h4 className="mb-1">Wallet Balance</h4>
          <TokenAmount amount={keepToken.value} withIcon withMetricSuffix />
        </section>
        <section className="tile wallet-page__overview__staked-tokens">
          <h4 className="mb-2">Tokens Staked</h4>
          <ProgressBar
            value={totalOwnedStakedBalance}
            total={totalBalance}
            color={colors.mint80}
            bgColor={colors.mint20}
          >
            <div className="circular-progress-bar-percentage-label-wrapper">
              <ProgressBar.Circular radius={82} barWidth={16} />
              <ProgressBar.PercentageLabel text="Staked" />
            </div>
            <ProgressBar.Legend
              leftValueLabel="Unstaked"
              valueLabel="Staked"
              displayLegendValuFn={displayAmountWithMetricSuffix}
            />
          </ProgressBar>
        </section>
        <section className="tile wallet-page__overview__stake-form">
          <header className="flex row center mb-1">
            <h3>Stake Wallet Tokens&nbsp;</h3>
            <ResourceTooltip {...resourceTooltipProps.delegation} />
          </header>
          <DelegateStakeForm
            onSubmit={onSubmitDelegateStakeForm}
            minStake={minimumStake}
            availableToStake={keepToken.value}
          />
        </section>
      </section>
      <DelegationOverview
        delegations={ownedDelegations}
        undelegations={ownedUndelegations}
        isFetching={isDelegationDataFetching}
        topUps={topUps}
        areTopUpsFetching={areTopUpsFetching}
        undelegationPeriod={undelegationPeriod}
        initializationPeriod={initializationPeriod}
        keepTokenBalance={keepToken.value}
      />
    </>
  )
}

const renderWalletTokensPageComponent = (onSubmitDelegateStakeForm) => (
  <WalletTokensPageComponent
    onSubmitDelegateStakeForm={onSubmitDelegateStakeForm}
  />
)

const WalletTokensPage = () => (
  <DelegationPageWrapper render={renderWalletTokensPageComponent} />
)

WalletTokensPage.route = {
  title: "Wallet Tokens",
  path: "/delegations/wallet",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export { WalletTokensPage }
