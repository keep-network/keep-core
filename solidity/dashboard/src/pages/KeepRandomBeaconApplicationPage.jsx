import React, { useCallback, useMemo, useState } from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import AuthorizationHistory from "../components/AuthorizationHistory"
import * as Icons from "../components/Icons"
import { useFetchData } from "../hooks/useFetchData"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { LoadingOverlay } from "../components/Loadable"
import { beaconAuthorizationService } from "../services/beacon-authorization.service"
import { isSameEthAddress } from "../utils/general.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import { authorizeOperatorContract } from "../actions/web3"
import { connect } from "react-redux"
import { getKeepRandomBeaconOperatorAddress } from "../contracts"

const KeepRandomBeaconApplicationPage = ({ authorizeOperatorContract }) => {
  const [selectedOperator, setOperator] = useState({})

  const [{ data, isFetching }, updateData] = useFetchData(
    beaconAuthorizationService.fetchRandomBeaconAuthorizationData,
    []
  )

  const onAuthorizationSuccessCallback = useCallback(
    (contractName, operatorAddress) => {
      const {
        indexInArray: operatorIndexInArray,
        obj: obsoleteOperator,
      } = findIndexAndObject(
        "operatorAddress",
        operatorAddress,
        data,
        compareEthAddresses
      )
      if (operatorIndexInArray === null) {
        return
      }
      const {
        indexInArray: contractIndexInArray,
        obj: obsoleteContract,
      } = findIndexAndObject(
        "contractName",
        contractName,
        obsoleteOperator.contracts
      )
      const updatedContracts = [...obsoleteOperator.contracts]
      updatedContracts[contractIndexInArray] = {
        ...obsoleteContract,
        isAuthorized: true,
      }
      const updatedOperators = [...data]
      updatedOperators[operatorIndexInArray] = {
        ...obsoleteOperator,
        contracts: updatedContracts,
      }

      updateData(updatedOperators)
    },
    [data, updateData]
  )

  const authorizeContract = useCallback(
    async (data, awaitingPromise) => {
      const { operatorAddress } = data
      const operatorContractAddress = getKeepRandomBeaconOperatorAddress()

      authorizeOperatorContract(
        { operatorAddress, operatorContractAddress },
        awaitingPromise
      )
    },
    [authorizeOperatorContract]
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
    <PageWrapper
      className=""
      title="Random Beacon"
      nextPageLink="/rewards/random-beacon"
      nextPageTitle="Rewards"
      nextPageIcon={Icons.KeepBlackGreen}
    >
      <nav className="mb-2">
        <a
          href="https://keep.network/"
          className="arrow-link h4"
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
    </PageWrapper>
  )
}

const toAuthHistoryData = (authData) => ({
  ...authData,
  ...authData.contracts[0],
})

const mapDispatchToProps = {
  authorizeOperatorContract,
}

export default connect(
  null,
  mapDispatchToProps
)(KeepRandomBeaconApplicationPage)
