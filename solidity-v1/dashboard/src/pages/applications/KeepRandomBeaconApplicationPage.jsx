import React, { useEffect, useCallback, useMemo, useState } from "react"
import { useSelector, useDispatch } from "react-redux"
import AuthorizeContracts from "../../components/AuthorizeContracts"
import AuthorizationHistory from "../../components/AuthorizationHistory"
import { LoadingOverlay } from "../../components/Loadable"
import { isSameEthAddress } from "../../utils/general.utils"
import DataTableSkeleton from "../../components/skeletons/DataTableSkeleton"
import { authorizeOperatorContract } from "../../actions/web3"
import { getKeepRandomBeaconOperatorAddress } from "../../contracts"
import EmptyStatePage from "./EmptyStatePage"
import { useWeb3Address } from "../../components/WithWeb3Context"
import {
  FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
  KEEP_RANDOM_BEACON_AUTHORIZED,
} from "../../actions"
import { LINK } from "../../constants/constants"

const KeepRandomBeaconApplicationPage = () => {
  const dispatch = useDispatch()
  const address = useWeb3Address()
  const [selectedOperator, setOperator] = useState({})
  const { isFetching, authData: data } = useSelector(
    (state) => state.authorization
  )

  useEffect(() => {
    dispatch({
      type: FETCH_KEEP_RANDOM_BEACON_AUTH_DATA_REQUEST,
      payload: { address },
    })
  }, [dispatch, address])

  const onAuthorizationSuccessCallback = useCallback(
    (contractName, operatorAddress) => {
      dispatch({
        type: KEEP_RANDOM_BEACON_AUTHORIZED,
        payload: { contractName, operatorAddress },
      })
    },
    [dispatch]
  )

  const authorizeContract = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress } = data
      const operatorContractAddress = getKeepRandomBeaconOperatorAddress()

      dispatch(
        authorizeOperatorContract(
          { operatorAddress, operatorContractAddress },
          awaitingPromise
        )
      )
    },
    [dispatch]
  )

  const authorizeContractsData = useMemo(() => {
    if (!selectedOperator.operatorAddress)
      return data.filter((authData) => !authData.contracts[0].isAuthorized)
    return data.filter(
      ({ operatorAddress, contracts }) =>
        !contracts[0].isAuthorized &&
        isSameEthAddress(operatorAddress, selectedOperator.operatorAddress)
    )
  }, [data, selectedOperator.operatorAddress])

  const authorizationHistoryData = useMemo(() => {
    if (!selectedOperator.operatorAddress)
      return data
        .filter((authData) => authData.contracts[0].isAuthorized)
        .map(toAuthHistoryData)
    return data
      .filter(
        ({ operatorAddress, contracts }) =>
          contracts[0].isAuthorized &&
          isSameEthAddress(operatorAddress, selectedOperator.operatorAddress)
      )
      .map(toAuthHistoryData)
  }, [data, selectedOperator.operatorAddress])

  return (
    <>
      <nav className="mb-2">
        <a
          href={LINK.keepWebsite}
          className="h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          Keep Website
        </a>
      </nav>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={
          <DataTableSkeleton columns={4} subtitleWidth="30%" />
        }
      >
        <AuthorizeContracts
          filterDropdownOptions={data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={authorizeContractsData}
          onAuthorizeBtn={authorizeContract}
          onSuccessCallback={onAuthorizationSuccessCallback}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={isFetching}
        skeletonComponent={<DataTableSkeleton columns={4} subtitleWidth="0" />}
      >
        <AuthorizationHistory contracts={authorizationHistoryData} />
      </LoadingOverlay>
    </>
  )
}

const toAuthHistoryData = (authData) => ({
  ...authData,
  ...authData.contracts[0],
})

KeepRandomBeaconApplicationPage.route = {
  title: "Keep Random Beacon",
  path: "/applications/random-beacon",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

export default KeepRandomBeaconApplicationPage
