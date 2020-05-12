import React, { useCallback } from "react"
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

const TBTCApplicationPage = () => {
  // fetch data from service
  const initialData = []
  const [state, , refreshData] = useFetchData(
    tbtcAuthorizationService.fetchTBTCAuthorizationData,
    initialData
  )
  // fetch bonding data
  const [bondingState, updateBondinData] = useFetchData(
    tbtcAuthorizationService.getBondingData,
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

  const web3Context = useWeb3Context()
  const showMessage = useShowMessage()

  const authorizeContract = useCallback(
    async (data, transactionHashCallback) => {
      const { operatorAddress, contractName } = data
      const serviceMethod =
        contractName === "TBTCSystem"
          ? tbtcAuthorizationService.authorizeTBTCSystem
          : tbtcAuthorizationService.authorizeBondedECDSAKeepFactory
      try {
        await serviceMethod.authorizeTBTCSystem(
          web3Context,
          operatorAddress,
          transactionHashCallback
        )
        showMessage({
          type: messageType.SUCCESS,
          title: "Success",
          content: "Authorization transaction successfully completed",
        })
      } catch (error) {
        showMessage({
          type: messageType.ERROR,
          title: "Authorization action has failed ",
          content: error.message,
        })
        throw error
      }
    },
    [showMessage, web3Context]
  )

  return (
    <PageWrapper
      className=""
      title="tBTC"
      nextPageLink="/rewards"
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
      <AuthorizeContracts
        data={state.data}
        onAuthorizeBtn={authorizeContract}
        onAuthorizeSuccessCallback={refreshData}
      />
      <BondingSection data={bondingState.data} />
    </PageWrapper>
  )
}

export default TBTCApplicationPage
