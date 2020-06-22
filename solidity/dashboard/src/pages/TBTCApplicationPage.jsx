import React, { useCallback, useState, useMemo } from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
// import * as Icons from "../components/Icons"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useFetchData } from "../hooks/useFetchData"
import { BondingSection } from "../components/BondingSection"
import { useWeb3Context } from "../components/WithWeb3Context"
import { useShowMessage, messageType } from "../components/Message"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { add, sub } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import { KEEP_BONDING_CONTRACT_NAME } from "../constants/constants"
import { LoadingOverlay } from "../components/Loadable"
import { isSameEthAddress } from "../utils/general.utils"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"

const initialData = []
const TBTCApplicationPage = () => {
  const web3Context = useWeb3Context()
  const showMessage = useShowMessage()
  const [selectedOperator, setOperator] = useState({})

  // fetch data from service
  const [tbtcAuthState, updateTbtcAuthData] = useFetchData(
    tbtcAuthorizationService.fetchTBTCAuthorizationData,
    initialData
  )
  // fetch bonding data
  const [bondingState, updateBondinData] = useFetchData(
    tbtcAuthorizationService.fetchBondingData,
    initialData
  )

  const unbondedValueUpdated = useCallback(
    (event, arithmeticOpration = add) => {
      const {
        returnValues: { operator, amount },
      } = event
      const { indexInArray, obj: obsoleteData } = findIndexAndObject(
        "operatorAddress",
        operator,
        bondingState.data,
        compareEthAddresses
      )
      if (indexInArray === null) {
        return
      }

      const availableETHInWei = arithmeticOpration(
        obsoleteData.availableETHInWei,
        amount
      ).toString()
      const availableETH = web3Utils.fromWei(availableETHInWei, "ether")
      const updatedBondinData = [...bondingState.data]
      updatedBondinData[indexInArray] = {
        ...obsoleteData,
        availableETH,
        availableETHInWei,
      }
      updateBondinData(updatedBondinData)
    },
    [updateBondinData, bondingState.data]
  )

  const subscribeToUnbondedValueDepositedCallback = (event) => {
    unbondedValueUpdated(event)
  }

  useSubscribeToContractEvent(
    KEEP_BONDING_CONTRACT_NAME,
    "UnbondedValueDeposited",
    subscribeToUnbondedValueDepositedCallback
  )

  const unbondedValueWithdrawnCallback = (event) => {
    unbondedValueUpdated(event, sub)
  }

  useSubscribeToContractEvent(
    KEEP_BONDING_CONTRACT_NAME,
    "UnbondedValueWithdrawn",
    unbondedValueWithdrawnCallback
  )

  const onSuccessCallback = useCallback(
    (contractName, operatorAddress, isAuthorized = true) => {
      const {
        indexInArray: operatorIndexInArray,
        obj: obsoleteOperator,
      } = findIndexAndObject(
        "operatorAddress",
        operatorAddress,
        tbtcAuthState.data,
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
        isAuthorized,
      }
      const updatedOperators = [...tbtcAuthState.data]
      updatedOperators[operatorIndexInArray] = {
        ...obsoleteOperator,
        contracts: updatedContracts,
      }

      updateTbtcAuthData(updatedOperators)
    },
    [updateTbtcAuthData, tbtcAuthState.data]
  )

  const authorizeContract = useCallback(
    async (data, transactionHashCallback) => {
      const { operatorAddress, contractName } = data
      const serviceMethod =
        contractName === "TBTCSystem"
          ? tbtcAuthorizationService.authorizeTBTCSystem
          : tbtcAuthorizationService.authorizeBondedECDSAKeepFactory
      try {
        await serviceMethod(
          web3Context,
          operatorAddress,
          transactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Authorization successfully completed",
        })
        setTimeout(() => onSuccessCallback(contractName, operatorAddress), 5000)
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Authorization has failed",
          content: error.message,
        })
        throw error
      }
    },
    [showMessage, web3Context, onSuccessCallback]
  )

  const deauthorizeTBTCSystem = useCallback(
    async (data, transactionHashCallback) => {
      const { operatorAddress } = data
      try {
        await tbtcAuthorizationService.deauthorizeTBTCSystem(
          web3Context,
          operatorAddress,
          transactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Deauthorization successfully completed",
        })
        setTimeout(
          () => onSuccessCallback("TBTCSystem", operatorAddress, false),
          5000
        )
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Deauthorization has failed",
          content: error.message,
        })
        throw error
      }
    },
    [showMessage, web3Context, onSuccessCallback]
  )

  const tbtcAuthData = useMemo(() => {
    if (!selectedOperator.operatorAddress) {
      return tbtcAuthState.data
    }
    return tbtcAuthState.data.filter((data) =>
      isSameEthAddress(data.operatorAddress, selectedOperator.operatorAddress)
    )
  }, [selectedOperator.operatorAddress, tbtcAuthState.data])

  return (
    <PageWrapper
      className=""
      title="tBTC"
      // The rewards page for the tbtc is not yet implemented
      // nextPageLink="/rewards"
      // nextPageTitle="Rewards"
      // nextPageIcon={Icons.TBTC}
    >
      <nav className="mb-2">
        <a
          href="https://tbtc.network/"
          className="arrow-link h4"
          rel="noopener noreferrer"
          target="_blank"
        >
          tBTC Website
        </a>
      </nav>
      <LoadingOverlay
        isFetching={tbtcAuthState.isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <AuthorizeContracts
          filterDropdownOptions={tbtcAuthState.data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={tbtcAuthData}
          onAuthorizeBtn={authorizeContract}
          onDeauthorizeBtn={deauthorizeTBTCSystem}
        />
      </LoadingOverlay>
      <LoadingOverlay
        isFetching={bondingState.isFetching}
        skeletonComponent={<DataTableSkeleton />}
      >
        <BondingSection data={bondingState.data} />
      </LoadingOverlay>
    </PageWrapper>
  )
}

export default TBTCApplicationPage
