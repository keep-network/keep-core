import React, { useCallback, useMemo, useState } from "react"
import { useSelector } from "react-redux"
import EmptyStatePage from "./EmptyStatePage"
import { isSameEthAddress } from "../../utils/general.utils"
import { LoadingOverlay } from "../../components/Loadable"
import DataTableSkeleton from "../../components/skeletons/DataTableSkeleton"
import AuthorizeThresholdContracts from "../../components/threshold/AuthorizeThresholdContracts"
import ThresholdAuthorizationHistory from "../../components/threshold/ThresholdStakingAuthorizationHistory"
import { MODAL_TYPES } from "../../constants/constants"
import { useModal } from "../../hooks/useModal"
import AuthorizeStakesBanner from "../../components/threshold/AuthorizeStakesBanner"
import Tile from "../../components/Tile"
import OnlyIf from "../../components/OnlyIf"

const ThresholdApplicationPage = () => {
  const [selectedOperator, setOperator] = useState({})
  const { openModal } = useModal()
  const thresholdAuthState = useSelector(
    (state) => state.thresholdAuthorization
  )

  const authorizeContract = useCallback(
    async (data) => {
      const {
        owner,
        operatorAddress,
        authorizerAddress,
        beneficiaryAddress,
        stakeAmount,
        isFromGrant,
        canBeMovedToT,
      } = data

      if (isFromGrant && !canBeMovedToT) {
        openModal(MODAL_TYPES.ContactYourGrantManagerWarning)
      } else {
        openModal(MODAL_TYPES.AuthorizeAndStakeOnThreshold, {
          keepAmount: stakeAmount,
          owner: owner,
          operator: operatorAddress,
          beneficiary: beneficiaryAddress,
          authorizer: authorizerAddress,
          isAuthorized: false,
        })
      }
    },
    [openModal]
  )

  const stakeToT = useCallback(
    async (data) => {
      const {
        owner,
        operatorAddress,
        authorizerAddress,
        beneficiaryAddress,
        stakeAmount,
        isFromGrant,
        canBeMovedToT,
      } = data

      if (isFromGrant && !canBeMovedToT) {
        openModal(MODAL_TYPES.ContactYourGrantManagerWarning)
      } else {
        openModal(MODAL_TYPES.StakeOnThresholdWithoutAuthorization, {
          keepAmount: stakeAmount,
          owner: owner,
          operator: operatorAddress,
          beneficiary: beneficiaryAddress,
          authorizer: authorizerAddress,
          isAuthorized: true,
        })
      }
    },
    [openModal]
  )

  const stakesToAuthOrMoveToT = useMemo(() => {
    const stakesNotMovedToT = thresholdAuthState.authData.filter((dataObj) => {
      return !dataObj.isStakedToT
    })
    if (!selectedOperator.operatorAddress) {
      return stakesNotMovedToT
    }
    return stakesNotMovedToT.filter((data) =>
      isSameEthAddress(data.operatorAddress, selectedOperator.operatorAddress)
    )
  }, [selectedOperator.operatorAddress, thresholdAuthState.authData])

  const authorizationHistoryData = useMemo(() => {
    if (!selectedOperator.operatorAddress)
      return thresholdAuthState.authData
        .filter((authData) => authData.contract.isAuthorized)
        .map(toAuthHistoryData)
    return thresholdAuthState.authData
      .filter(
        ({ operatorAddress, contract }) =>
          contract.isAuthorized &&
          isSameEthAddress(operatorAddress, selectedOperator.operatorAddress)
      )
      .map(toAuthHistoryData)
  }, [thresholdAuthState.authData, selectedOperator.operatorAddress])

  return (
    <>
      <LoadingOverlay
        isFetching={thresholdAuthState.isFetching}
        skeletonComponent={
          <DataTableSkeleton columns={4} subtitleWidth="40%" />
        }
      >
        <AuthorizeStakesBanner stakesToAuthOrMoveToT={stakesToAuthOrMoveToT} />
        <AuthorizeThresholdContracts
          filterDropdownOptions={thresholdAuthState.authData}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={stakesToAuthOrMoveToT}
          onAuthorizeBtn={authorizeContract}
          onStakeBtn={stakeToT}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={thresholdAuthState.isFetching}
        skeletonComponent={<DataTableSkeleton columns={4} subtitleWidth="0" />}
      >
        <OnlyIf condition={thresholdAuthState.authData.length > 0}>
          <ThresholdAuthorizationHistory contracts={authorizationHistoryData} />
        </OnlyIf>
        <OnlyIf condition={thresholdAuthState.authData.length === 0}>
          <Tile className={"tile threshold-staking__no-data-tile"}>
            <div className={"text-center"}>
              <h3 className={"threshold-staking__title text-grey-60 mb-1"}>
                Threshold Staking
              </h3>
              <span className={"text-grey-60"}>
                Authorize the staking contract{<br />}above to stake and earn
                rewards.
              </span>
            </div>
          </Tile>
        </OnlyIf>
      </LoadingOverlay>
    </>
  )
}

const toAuthHistoryData = (authData) => ({
  ...authData,
  ...authData.contract,
})

ThresholdApplicationPage.route = {
  title: "Threshold",
  path: "/applications/threshold",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
  withNewLabel: true,
}

export default ThresholdApplicationPage
