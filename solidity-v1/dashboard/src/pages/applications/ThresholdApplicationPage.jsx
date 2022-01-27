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

const initialData = []
const ThresholdApplicationPage = ({ authorizeOperatorContract }) => {
  const [selectedOperator, setOperator] = useState({})
  const address = useWeb3Address()

  const [
    thresholdAuthState,
    updateThresholdAuthData,
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
      console.log("operatorAddress", operatorAddress)
      console.log("operatorContractAddress", operatorContractAddress)
      authorizeOperatorContract(
        { operatorAddress, operatorContractAddress },
        awaitingPromise
      )
    },
    [authorizeOperatorContract]
  )

  const stakeToT = () => console.log(updateThresholdAuthData)

  // const stakeToT = useCallback(
  //   async (data, awaitingPromise) => {},
  //   [authorizeOperatorContract]
  // )

  const thresholdAuthData = useMemo(() => {
    if (!selectedOperator.operatorAddress) {
      return thresholdAuthState.data
    }
    return thresholdAuthState.data.filter((data) =>
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

  console.log("auth history data", authorizationHistoryData)
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
          // onSuccessCallback={onSuccessCallback}
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
