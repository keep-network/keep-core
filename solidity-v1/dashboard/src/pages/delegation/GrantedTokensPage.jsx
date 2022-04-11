import React, { useEffect, useCallback, useMemo } from "react"
import { useLocation } from "react-router-dom"
import { useSelector } from "react-redux"
import EmptyStateComponent from "./EmptyStatePage"
import DelegateStakeForm from "../../components/DelegateStakeForm"
import {
  TokenGrantDetails,
  TokenGrantStakedDetails,
  TokenGrantUnlockingdDetails,
  TokenGrantWithdrawnTokensDetails,
} from "../../components/TokenGrantOverview"
import { DelegationPageWrapper } from "./index"
import DelegationOverview from "../../components/DelegationOverview"
import ResourceTooltip from "../../components/ResourceTooltip"
import { useState } from "react"
import { isEmptyArray } from "../../utils/array.utils"
import { isEmptyObj } from "../../utils/general.utils"
import { add } from "../../utils/arithmetics.utils"
import { usePrevious } from "../../hooks/usePrevious"
import { CompoundDropdown as Dropdown } from "../../components/Dropdown"
import * as Icons from "../../components/Icons"
import useReleaseTokens from "../../hooks/useReleaseTokens"
import resourceTooltipProps from "../../constants/tooltips"
import useDelegationsWithTAuthData from "../../hooks/useDelegationsWithTAuthData"

const filterBySelectedGrant = (selectedGrant) => (delegation) =>
  selectedGrant.id && delegation.grantId === selectedGrant.id

const GrantedTokensPageComponent = ({ onSubmitDelegateStakeForm }) => {
  const [selectedGrant, setSelectedGrant] = useState({})
  const previousGrant = usePrevious(selectedGrant)
  const releaseTokens = useReleaseTokens()
  const { hash } = useLocation()

  const {
    undelegations,
    undelegationPeriod,
    initializationPeriod,
    topUps,
    areTopUpsFetching,
    minimumStake,
  } = useSelector((state) => state.staking)

  const { grants, isFetching: areGrantsFetching } = useSelector(
    (state) => state.tokenGrants
  )

  const delegationsWithTAuthData = useDelegationsWithTAuthData()

  useEffect(() => {
    if (!isEmptyArray(grants) && isEmptyObj(selectedGrant)) {
      const lookupGrantId = hash.substr(1)
      const grantByHash = grants.find((grant) => grant.id === lookupGrantId)
      setSelectedGrant(grantByHash || grants[0])
    } else if (!isEmptyObj(selectedGrant) && isEmptyArray(grants)) {
      setSelectedGrant({})
    }
  }, [grants, selectedGrant, previousGrant.id, hash])

  const onSelectGrant = useCallback((selectedGrant) => {
    setSelectedGrant(selectedGrant)
  }, [])

  const grantDelegationsWithTAuthData = useMemo(() => {
    return delegationsWithTAuthData.filter(filterBySelectedGrant(selectedGrant))
  }, [delegationsWithTAuthData, selectedGrant])

  const grantUndelegations = useMemo(() => {
    return undelegations.filter(filterBySelectedGrant(selectedGrant))
  }, [undelegations, selectedGrant])

  const selectedGrantStakedAmount = useMemo(() => {
    if (!selectedGrant.id) return 0

    return [...grantDelegationsWithTAuthData, ...grantUndelegations]
      .filter((delegation) => delegation.grantId === selectedGrant.id)
      .map((grantDelegation) => grantDelegation.amount)
      .reduce(add, 0)
  }, [grantUndelegations, grantDelegationsWithTAuthData, selectedGrant.id])

  const onWithdrawTokens = useCallback(
    async (awaitingPromise) => {
      releaseTokens(selectedGrant, awaitingPromise)
    },
    [releaseTokens, selectedGrant]
  )

  const onSubmit = useCallback(
    (values, meta) => {
      values.grantData = selectedGrant
      onSubmitDelegateStakeForm(values, meta)
    },
    [onSubmitDelegateStakeForm, selectedGrant]
  )

  return (
    <>
      <section className="granted-page__overview-layout">
        <Dropdown
          selectedItem={selectedGrant}
          onSelect={onSelectGrant}
          comparePropertyName="id"
          className="granted-page__grants-dropdown"
          rounded
        >
          <Dropdown.Trigger>
            <div className="flex row center">
              <Icons.Grant width={14} height={14} />
              <span className="text-instruction text-grey-60 ml-1">
                Switch Grant
              </span>
            </div>
          </Dropdown.Trigger>
          <Dropdown.Options>{grants.map(renderGrant)}</Dropdown.Options>
        </Dropdown>
        <section className="tile granted-page__overview__grant-details">
          <h4 className="mb-1">Grant Allocation</h4>
          <TokenGrantDetails
            selectedGrant={selectedGrant}
            availableAmount={selectedGrant.availableToStake}
          />
        </section>
        <section className="tile granted-page__overview__staked-tokens">
          <h4 className="mb-2">Tokens Staked</h4>
          <TokenGrantStakedDetails
            selectedGrant={selectedGrant}
            stakedAmount={selectedGrantStakedAmount}
          />
        </section>
        <section className="tile granted-page__overview__stake-form">
          <header className="flex row center mb-1">
            <h3>Stake Granted Tokens&nbsp;</h3>
            <ResourceTooltip {...resourceTooltipProps.delegation} />
          </header>
          <DelegateStakeForm
            onSubmit={onSubmit}
            minStake={minimumStake}
            availableToStake={selectedGrant.availableToStake || 0}
          />
        </section>
        <section className="tile granted-page__overview__withdraw-tokens">
          <h4 className="mb-2">Withdraw Unlocked Tokens</h4>
          <TokenGrantWithdrawnTokensDetails
            selectedGrant={selectedGrant}
            onWithdrawnBtn={onWithdrawTokens}
          />
        </section>
        <section className="tile granted-page__overview__unlocked-tokens">
          <h4 className="mb-2">Tokens Unlocking Progress</h4>
          <TokenGrantUnlockingdDetails selectedGrant={selectedGrant} />
        </section>
      </section>
      <DelegationOverview
        delegationsWithTAuthData={grantDelegationsWithTAuthData}
        undelegations={grantUndelegations}
        isFetching={areGrantsFetching}
        topUps={topUps}
        areTopUpsFetching={areTopUpsFetching}
        undelegationPeriod={undelegationPeriod}
        initializationPeriod={initializationPeriod}
        selectedGrant={selectedGrant}
        grants={grants}
        context="granted"
      />
    </>
  )
}

const renderGrant = (grant) => (
  <Dropdown.Option key={grant.id} value={grant}>
    Grant #{grant.id}
  </Dropdown.Option>
)

const renderGrantedTokensPageComponent = (onSubmitDelegateStakeForm) => (
  <GrantedTokensPageComponent
    onSubmitDelegateStakeForm={onSubmitDelegateStakeForm}
  />
)

const GrantedTokensPage = () => (
  <DelegationPageWrapper render={renderGrantedTokensPageComponent} />
)

GrantedTokensPage.route = {
  title: "Granted Tokens",
  path: "/delegations/granted",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export { GrantedTokensPage }
