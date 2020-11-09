import React, { useEffect, useCallback, useMemo } from "react"
import { useSelector, useDispatch } from "react-redux"
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
import { useState } from "react"
import { isEmptyArray } from "../../utils/array.utils"
import { isEmptyObj } from "../../utils/general.utils"
import { usePrevious } from "../../hooks/usePrevious"
import { CompoundDropdown as Dropdown } from "../../components/Dropdown"
import * as Icons from "../../components/Icons"
import { ContractsLoaded } from "../../contracts"
import { useModal } from "../../hooks/useModal"
import { ViewAddressInBlockExplorer } from "../../components/ViewInBlockExplorer"
import { releaseTokens } from "../../actions/web3"
import { withConfirmationModal } from "../../components/ConfirmationModal"

const filterBySelectedGrant = (selectedGrant) => (delegation) =>
  selectedGrant.id && delegation.grantId === selectedGrant.id

const GrantedTokensPageComponent = ({ onSubmitDelegateStakeForm }) => {
  const [selectedGrant, setSelectedGrant] = useState({})
  const previousGrant = usePrevious(selectedGrant)
  const dispatch = useDispatch()
  const { openConfirmationModal } = useModal()

  const {
    undelegations,
    delegations,
    undelegationPeriod,
    initializationPeriod,
    topUps,
    areTopUpsFetching,
    minimumStake,
  } = useSelector((state) => state.staking)

  const { grants, isFetching: areGrantsFetching } = useSelector(
    (state) => state.tokenGrants
  )

  useEffect(() => {
    if (!isEmptyArray(grants) && isEmptyObj(selectedGrant)) {
      setSelectedGrant(grants[0])
    }
  }, [grants, selectedGrant, previousGrant.id])

  const onSelectGrant = useCallback((selectedGrant) => {
    setSelectedGrant(selectedGrant)
  }, [])

  const grantDelegations = useMemo(() => {
    return delegations.filter(filterBySelectedGrant(selectedGrant))
  }, [delegations, selectedGrant])

  const grantUndelegations = useMemo(() => {
    return undelegations.filter(filterBySelectedGrant(selectedGrant))
  }, [undelegations, selectedGrant])

  const onWithdrawTokens = useCallback(
    async (awaitingPromise) => {
      const { escrowOperatorsToWithdraw } = selectedGrant

      if (!isEmptyArray(escrowOperatorsToWithdraw)) {
        const { tokenStakingEscrow } = await ContractsLoaded
        await openConfirmationModal(
          {
            modalOptions: { title: "Are you sure?" },
            title: "You’re about to release tokens.",
            escrowAddress: tokenStakingEscrow.options.address,
            btnText: "release",
            confirmationText: "RELEASE",
          },
          withConfirmationModal(ConfirmWithdrawModal)
        )
      }
      dispatch(releaseTokens(selectedGrant, awaitingPromise))
    },
    [dispatch, openConfirmationModal, selectedGrant]
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
            availableAmount={0}
          />
        </section>
        <section className="tile granted-page__overview__staked-tokens">
          <h4 className="mb-2">Tokens Staked</h4>
          <TokenGrantStakedDetails
            selectedGrant={selectedGrant}
            stakedAmount={0}
          />
        </section>
        <section className="tile granted-page__overview__stake-form">
          <h3 className="mb-1">Stake Granted Tokens</h3>
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
        delegations={grantDelegations}
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
    Grant {grant.id}
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
  path: "/delegation/grant",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStateComponent,
}

export { GrantedTokensPage }

const ConfirmWithdrawModal = ({ escrowAddress }) => {
  return (
    <>
      <span>You have deposited tokens in the</span>&nbsp;
      <ViewAddressInBlockExplorer
        text="TokenStakingEscrow contract"
        address={escrowAddress}
      />
      <p>
        To withdraw all tokens it may be necessary to confirm more than one
        transaction.
      </p>
    </>
  )
}
