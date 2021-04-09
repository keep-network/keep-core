import React, { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import { useWeb3Context } from "../components/WithWeb3Context"
import { WALLETS } from "../constants/constants"
import { useModal } from "./useModal"
import { WalletSelectionModal } from "../components/WalletSelectionModal"
import { useLocation, useHistory } from "react-router-dom"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl"

const useSubscribeToConnectorEvents = () => {
  const dispatch = useDispatch()
  const { isConnected, connector, yourAddress, web3 } = useWeb3Context()
  const { openModal } = useModal()
  const { transactionQueue } = useSelector((state) => state.transactions)
  const history = useHistory()
  const location = useLocation()
  const walletAddressFromUrl = useWalletAddressFromUrl()

  useEffect(() => {
    const accountChangedHandler = (address) => {
      dispatch({ type: "app/account_changed", payload: { address } })
    }

    const disconnectHandler = () => {
      dispatch({ type: "app/logout" })
    }

    const showChooseWalletModal = (payload) => {
      dispatch({
        type: "transactions/transaction_added_to_queue",
        payload: payload,
      })
      openModal(<WalletSelectionModal />, {
        title: "Select Wallet",
      })
    }

    const executeTransactionsInQueue = async (transactions) => {
      if (transactions.length > 0) {
        dispatch({
          type: "transactions/clear_queue",
        })
        for (const transaction of transactions) {
          await web3.eth.currentProvider.sendAsync(transaction)
        }
      }
    }

    if (isConnected && connector) {
      dispatch({ type: "app/login", payload: { address: yourAddress } })
      if (connector.name === WALLETS.METAMASK.name) {
        connector.getProvider().on("accountsChanged", accountChangedHandler)
        connector.getProvider().on("chainChanged", disconnectHandler)
      }

      if (connector.name === WALLETS.EXPLORER_MODE.name) {
        connector.eventEmitter.on(
          "chooseWalletAndSendTransaction",
          showChooseWalletModal
        )
      } else {
        executeTransactionsInQueue(transactionQueue)
        if (walletAddressFromUrl) {
          const newPath = location.pathname.replace(
            "/" + walletAddressFromUrl,
            ""
          )
          history.push({ pathname: newPath })
        }
      }
    }

    return () => {
      if (connector && connector.name === WALLETS.METAMASK.name) {
        connector
          .getProvider()
          .removeListener("accountsChanged", accountChangedHandler)
        connector.getProvider().removeListener("disconnect", disconnectHandler)
      }
    }
  }, [isConnected, connector, dispatch, yourAddress, web3])
}

export default useSubscribeToConnectorEvents
