import React, { useCallback, useMemo, useState } from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import AuthorizationHistory from "../components/AuthorizationHistory"
import * as Icons from "../components/Icons"
import { useFetchData } from "../hooks/useFetchData"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { useShowMessage, messageType } from "../components/Message"
import { useWeb3Context } from "../components/WithWeb3Context"
import { LoadingOverlay } from "../components/Loadable"
import { beaconAuthorizationService } from "../services/beacon-authorization.service"
import { isSameEthAddress } from "../utils/general.utils"

const KeepRandomBeaconApplicationPage = () => {
  const web3Context = useWeb3Context()
  const showMessage = useShowMessage()
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
    async (data, transactionHashCallback) => {
      const { operatorAddress, contractName } = data
      try {
        // TODO call service here
        // await serviceMethod(
        //   web3Context,
        //   operatorAddress,
        //   transactionHashCallback
        // )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Authorization transaction successfully completed",
        })
        setTimeout(
          () => onAuthorizationSuccessCallback(contractName, operatorAddress),
          5000
        )
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Authorization action has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [showMessage, onAuthorizationSuccessCallback]
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
      nextPageLink="/rewards"
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
      <LoadingOverlay isFetching={isFetching}>
        <AuthorizeContracts
          filterDropdownOptions={data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={authorizeContractsData}
          onAuthorizeBtn={authorizeContract}
        />
        <AuthorizationHistory contracts={authorizationHistoryData} />
      </LoadingOverlay>
    </PageWrapper>
  )
}

const toAuthHistoryData = (authData) => ({
  ...authData,
  ...authData.contracts[0],
})

export default KeepRandomBeaconApplicationPage
