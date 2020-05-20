import React, { useCallback, useState, useMemo } from "react"
import PageWrapper from "../components/PageWrapper"
import AuthorizeContracts from "../components/AuthorizeContracts"
import * as Icons from "../components/Icons"
import { tbtcAuthorizationService } from "../services/tbtc-authorization.service"
import { useFetchData } from "../hooks/useFetchData"
import { BondingSection } from "../components/BondingSection"
import { useWeb3Context } from "../components/WithWeb3Context"
import { useShowMessage, messageType } from "../components/Message"
import { useSubscribeToContractEvent } from "../hooks/useSubscribeToContractEvent"
import { findIndexAndObject, compareEthAddresses } from "../utils/array.utils"
import { add } from "../utils/arithmetics.utils"
import web3Utils from "web3-utils"
import { KEEP_BONDING_CONTRACT_NAME } from "../constants/constants"
import { LoadingOverlay } from "../components/Loadable"
import { isSameEthAddress } from "../utils/general.utils"

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

  const subscribeToUnbondedValueDepositedCallback = (event) => {
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

    const availableETHInWei = add(
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
  }

  useSubscribeToContractEvent(
    KEEP_BONDING_CONTRACT_NAME,
    "UnbondedValueDeposited",
    subscribeToUnbondedValueDepositedCallback
  )

  const onAuthorizationSuccessCallback = useCallback(
    (contractName, operatorAddress) => {
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
        isAuthorized: true,
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
    [showMessage, web3Context, onAuthorizationSuccessCallback]
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
      nextPageLink="/rewards/tbtc"
      nextPageTitle="Rewards"
      nextPageIcon={Icons.TBTC}
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
      <LoadingOverlay isFetching={tbtcAuthState.isFetching}>
        <AuthorizeContracts
          filterDropdownOptions={tbtcAuthState.data}
          onSelectOperator={setOperator}
          selectedOperator={selectedOperator}
          data={tbtcAuthData}
          onAuthorizeBtn={authorizeContract}
        />
      </LoadingOverlay>
      <LoadingOverlay isFetching={bondingState.isFetching}>
        <BondingSection data={bondingState.data} />
      </LoadingOverlay>
    </PageWrapper>
  )
}

export default TBTCApplicationPage
