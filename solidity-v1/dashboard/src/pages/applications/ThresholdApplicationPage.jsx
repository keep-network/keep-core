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

const ThresholdApplicationPage = () => {
  const [selectedOperator, setOperator] = useState({})
  const { openModal } = useModal()
  const thresholdAuthState = useSelector(
    (state) => state.thresholdAuthorization
  )

  const authorizeContract = useCallback(
    async (data) => {
      const {
        operatorAddress,
        authorizerAddress,
        beneficiaryAddress,
        stakeAmount,
      } = data
      openModal(MODAL_TYPES.AuthorizeAndStakeOnThreshold, {
        keepAmount: stakeAmount,
        operator: operatorAddress,
        beneficiary: beneficiaryAddress,
        authorizer: authorizerAddress,
        isAuthorized: false,
      })
    },
    [openModal]
  )

  const stakeToT = useCallback(
    async (data) => {
      const {
        operatorAddress,
        authorizerAddress,
        beneficiaryAddress,
        stakeAmount,
      } = data
      openModal(MODAL_TYPES.StakeOnThresholdWithoutAuthorization, {
        keepAmount: stakeAmount,
        operator: operatorAddress,
        beneficiary: beneficiaryAddress,
        authorizer: authorizerAddress,
        isAuthorized: true,
      })
    },
    [openModal]
  )

  const stakesToAuthOrMoveToT = useMemo(() => {
    const unauthorizedStakes = thresholdAuthState.authData.filter((dataObj) => {
      return !dataObj.isStakedToT || !dataObj.contracts[0].isAuthorized
    })
    if (!selectedOperator.operatorAddress) {
      return unauthorizedStakes
    }
    return unauthorizedStakes.filter((data) =>
      isSameEthAddress(data.operatorAddress, selectedOperator.operatorAddress)
    )
  }, [selectedOperator.operatorAddress, thresholdAuthState.authData])

  const authorizationHistoryData = useMemo(() => {
    if (!selectedOperator.operatorAddress)
      return thresholdAuthState.authData
        .filter((authData) => authData.contracts[0].isAuthorized)
        .map(toAuthHistoryData)
    return thresholdAuthState.authData
      .filter(
        ({ operatorAddress, contracts }) =>
          contracts[0].isAuthorized &&
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
        <AuthorizeStakesBanner
          numberOfStakesToAuthorize={stakesToAuthOrMoveToT.length}
        />
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
        <ThresholdAuthorizationHistory contracts={authorizationHistoryData} />
      </LoadingOverlay>
    </>
  )
}

const toAuthHistoryData = (authData) => ({
  ...authData,
  ...authData.contracts[0],
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
