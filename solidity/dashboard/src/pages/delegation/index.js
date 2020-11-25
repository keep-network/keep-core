import React, { useEffect, useCallback } from "react"
import { connect } from "react-redux"
import moment from "moment"
import { FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST } from "../../actions"
import { isEmptyArray } from "../../utils/array.utils"
import Banner, { BANNER_TYPE } from "../../components/Banner"
import Button from "../../components/Button"
import { useModal } from "../../hooks/useModal"
import CopyStakePage from "../CopyStakePage"
import PageWrapper from "../../components/PageWrapper"

import { WalletTokensPage } from "./WalletTokensPage"
import { GrantedTokensPage } from "./GrantedTokensPage"

const DelegationPage = ({ title, routes }) => {
  return <PageWrapper title={title} routes={routes} />
}

const DelegationPageWrapperComponent = ({
  fetchOldDelegations,
  oldDelegations,
  fetchGrants,
  fetchDelegations,
  fetchTopUps,
  delegateStake,
  initializationPeriod,
  children,
  ...restProps
}) => {
  useEffect(() => {
    fetchOldDelegations()
  }, [fetchOldDelegations])

  useEffect(() => {
    fetchGrants()
  }, [fetchGrants])

  useEffect(() => {
    fetchDelegations()
  }, [fetchDelegations])

  useEffect(() => {
    fetchTopUps()
  }, [fetchTopUps])

  const { openModal, openConfirmationModal } = useModal()

  const onSubmitDelegateStakeForm = useCallback(
    async (values, meta) => {
      await openConfirmationModal(
        confirmationModalOptions(initializationPeriod)
      )
      const grantData = values.grantData
        ? { ...values.grantData, grantId: values.grantData.id }
        : {}

      delegateStake(
        {
          ...values,
          ...grantData,
          amount: values.stakeTokens,
        },
        meta
      )
    },
    [delegateStake, initializationPeriod, openConfirmationModal]
  )

  return (
    <>
      {!isEmptyArray(oldDelegations) && (
        <Banner
          type={BANNER_TYPE.NOTIFICATION}
          withIcon
          title="New upgrade available for your stake delegations!"
          titleClassName="h4"
          subtitle="Upgrade now to keep earning rewards on your stake."
        >
          <Button
            className="btn btn-tertiary btn-sm ml-a"
            onClick={() => openModal(<CopyStakePage />, { isFullScreen: true })}
          >
            upgrade my stake
          </Button>
        </Banner>
      )}
      {restProps.render(onSubmitDelegateStakeForm)}
    </>
  )
}

const mapStateToProps = ({ copyStake, staking }) => {
  const { oldDelegations } = copyStake
  const { initializationPeriod } = staking

  return { oldDelegations, initializationPeriod }
}

const mapDispatchToProps = (dispatch) => {
  return {
    delegateStake: (values, meta) =>
      dispatch({
        type: "staking/delegate_request",
        payload: values,
        meta,
      }),
    fetchOldDelegations: () =>
      dispatch({ type: FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST }),
    fetchGrants: () => dispatch({ type: "token-grant/fetch_grants_request" }),
    fetchDelegations: () =>
      dispatch({ type: "staking/fetch_delegations_request" }),
    fetchTopUps: () => dispatch({ type: "staking/fetch_top_ups_request" }),
  }
}

const confirmationModalOptions = (initializationPeriod) => ({
  modalOptions: { title: "Initiate Delegation" },
  title: "You’re about to delegate stake.",
  subtitle: `You’re delegating KEEP tokens. You will be able to cancel the delegation for up to ${moment()
    .add(initializationPeriod, "seconds")
    .fromNow(true)}. After that time, you can undelegate your stake.`,
  btnText: "delegate",
  confirmationText: "DELEGATE",
})

export const DelegationPageWrapper = connect(
  mapStateToProps,
  mapDispatchToProps
)(DelegationPageWrapperComponent)

DelegationPage.route = {
  title: "Delegation",
  path: "/delegation",
  pages: [WalletTokensPage, GrantedTokensPage],
}

export default DelegationPage
