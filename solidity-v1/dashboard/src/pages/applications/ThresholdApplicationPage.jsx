import React, { useCallback, useEffect, useMemo, useState } from "react"
import { authorizeOperatorContract } from "../../actions/web3"
import { connect } from "react-redux"
import EmptyStatePage from "./EmptyStatePage"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { useFetchData } from "../../hooks/useFetchData"
import { getThresholdTokenStakingAddress } from "../../contracts"
import { thresholdAuthorizationService } from "../../services/threshold-authorization.service"
import { isSameEthAddress } from "../../utils/general.utils"
import { LoadingOverlay } from "../../components/Loadable"
import DataTableSkeleton from "../../components/skeletons/DataTableSkeleton"
import AuthorizeThresholdContracts from "../../components/threshold/AuthorizeThresholdContracts"
import ThresholdAuthorizationHistory from "../../components/threshold/ThresholdStakingAuthorizationHistory"
import { stakeKeepToT } from "../../actions/keep-to-t-staking"

const initialData = []
const ThresholdApplicationPage = ({
  authorizeOperatorContract,
  stakeKeepToT,
}) => {
  const [selectedOperator, setOperator] = useState({})
  const address = useWeb3Address()

  const [
    thresholdAuthState,
    updateThresholdAuthData,
    ,
    setThresholdAuthDataArgs,
  ] = useFetchData(
    thresholdAuthorizationService.fetchThresholdAuthorizationData,
    initialData,
    address
  )

  useEffect(() => {
    setThresholdAuthDataArgs([address])
  }, [setThresholdAuthDataArgs, address])

  const authorizeContract = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress } = data
      const operatorContractAddress = getThresholdTokenStakingAddress()
      authorizeOperatorContract(
        { operatorAddress, operatorContractAddress },
        awaitingPromise
      )
    },
    [authorizeOperatorContract]
  )

  const stakeToT = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress } = data
      stakeKeepToT(operatorAddress, awaitingPromise)
    },
    [stakeKeepToT]
  )

  const onSuccessCallback = () => console.log(updateThresholdAuthData)

  const thresholdAuthData = useMemo(() => {
    const thresholdData = thresholdAuthState.data.filter((dataObj) => {
      return !dataObj.isStakedToT || !dataObj.contracts[0].isAuthorized
    })
    if (!selectedOperator.operatorAddress) {
      return thresholdData
    }
    return thresholdData.filter((data) =>
      isSameEthAddress(data.operatorAddress, selectedOperator.operatorAddress)
    )
  }, [selectedOperator.operatorAddress, thresholdAuthState.data])

  const authorizationHistoryData = useMemo(() => {
    if (!selectedOperator.operatorAddress)
      return thresholdAuthState.data
        .filter((authData) => authData.contracts[0].isAuthorized)
        .map(toAuthHistoryData)
    return thresholdAuthState.data
      .filter(
        ({ operatorAddress, contracts }) =>
          contracts[0].isAuthorized &&
          isSameEthAddress(operatorAddress, selectedOperator.operatorAddress)
      )
      .map(toAuthHistoryData)
  }, [thresholdAuthState.data, selectedOperator.operatorAddress])

  return (
    <>
      <LoadingOverlay
        isFetching={thresholdAuthState.isFetching}
        skeletonComponent={
          <DataTableSkeleton columns={4} subtitleWidth="40%" />
        }
      >
        <AuthorizeThresholdContracts
          filterDropdownOptions={thresholdAuthState.data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={thresholdAuthData}
          onAuthorizeBtn={authorizeContract}
          onStakeBtn={stakeToT}
          onSuccessCallback={onSuccessCallback}
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

const mapDispatchToProps = {
  authorizeOperatorContract,
  stakeKeepToT,
}

const ConnectedThresholdApplicationPage = connect(
  null,
  mapDispatchToProps
)(ThresholdApplicationPage)

ConnectedThresholdApplicationPage.route = {
  title: "Threshold",
  path: "/applications/threshold",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default ConnectedThresholdApplicationPage
