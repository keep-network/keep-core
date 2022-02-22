import React, { useEffect, useCallback } from "react"
import { connect } from "react-redux"
import { FETCH_DELEGATIONS_FROM_OLD_STAKING_CONTRACT_REQUEST } from "../../actions"
import { isEmptyArray } from "../../utils/array.utils"
import Banner from "../../components/Banner"
import Button from "../../components/Button"
import { useModal } from "../../hooks/useModal"
import PageWrapper from "../../components/PageWrapper"
import * as Icons from "../../components/Icons"
import { WalletTokensPage } from "./WalletTokensPage"
import { GrantedTokensPage } from "./GrantedTokensPage"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { isDelegationExists } from "../../services/token-staking.service"
import { MODAL_TYPES } from "../../constants/constants"

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
  const yourAddress = useWeb3Address()

  useEffect(() => {
    fetchOldDelegations()
  }, [fetchOldDelegations, yourAddress])

  useEffect(() => {
    fetchGrants(yourAddress)
  }, [fetchGrants, yourAddress])

  useEffect(() => {
    fetchDelegations(yourAddress)
  }, [fetchDelegations, yourAddress])

  useEffect(() => {
    fetchTopUps(yourAddress)
  }, [fetchTopUps, yourAddress])

  const { openModal, openConfirmationModal } = useModal()

  const onSubmitDelegateStakeForm = useCallback(
    async (values, awaitingPromise) => {
      const { operatorAddress } = values
      try {
        if (await isDelegationExists(operatorAddress)) {
          openModal(MODAL_TYPES.DelegationAlreadyExists, { operatorAddress })
          throw new Error("Delegation already exists")
        }
        await openConfirmationModal(MODAL_TYPES.ConfirmDelegation, {
          initializationPeriod,
          isFromGrant: !!values.grantData,
        })
        const grantData = values.grantData
          ? { ...values.grantData, grantId: values.grantData.id }
          : {}

        delegateStake(
          {
            ...values,
            ...grantData,
            amount: values.stakeTokens,
          },
          awaitingPromise
        )
      } catch (error) {
        awaitingPromise.reject(error)
      }
    },
    [delegateStake, initializationPeriod, openConfirmationModal, openModal]
  )

  return (
    <>
      {!isEmptyArray(oldDelegations) && (
        <Banner className="banner--upgrade">
          <div className="flex row">
            <Banner.Icon icon={Icons.Alert} className="mr-1" />
            <div>
              <Banner.Title className="text-white h4">
                New upgrade available for your stake delegations!
              </Banner.Title>
              <Banner.Description className="text-grey-20">
                Upgrade now to keep earning rewards on your stake.
              </Banner.Description>
            </div>
            <Button
              className="btn btn-tertiary btn-sm ml-a"
              onClick={() => openModal(MODAL_TYPES.CopyStake)}
            >
              upgrade my stake
            </Button>
          </div>
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
    fetchGrants: (address) =>
      dispatch({
        type: "token-grant/fetch_grants_request",
        payload: { address },
      }),
    fetchDelegations: (address) =>
      dispatch({
        type: "staking/fetch_delegations_request",
        payload: { address },
      }),
    fetchTopUps: (address) =>
      dispatch({ type: "staking/fetch_top_ups_request", payload: { address } }),
  }
}

export const DelegationPageWrapper = connect(
  mapStateToProps,
  mapDispatchToProps
)(DelegationPageWrapperComponent)

DelegationPage.route = {
  title: "Delegations",
  path: "/delegations",
  pages: [GrantedTokensPage, WalletTokensPage],
}

export default DelegationPage
