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

      if (walletAddressFromUrl) {
        const newPathname = location.pathname.replace(
          "/" + walletAddressFromUrl,
          ""
        )
        history.push({ pathname: newPathname })
      }
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

    const sendTransactionsFromQueue = (transactions) => {
      const transactionObjects = transactions.map(
        (transaction) => transaction.params[0]
      )

      dispatch({
        type: "web3/send_raw_transaction_in_sequence",
        payload: transactionObjects,
      })
    }

    const executeTransactionsInQueue = (transactions) => {
      if (transactions.length > 0) {
        dispatch({
          type: "transactions/clear_queue",
        })
        sendTransactionsFromQueue(transactions)
      }
    }

    if (isConnected && connector) {
      dispatch({ type: "app/login", payload: { address: yourAddress } })
      connector.on("accountsChanged", accountChangedHandler)
      connector.once("disconnect", disconnectHandler)
      connector.on("chooseWalletAndSendTransaction", showChooseWalletModal)

      if (connector.name !== WALLETS.EXPLORER_MODE.name) {
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
      if (connector) {
        connector.removeListener("accountsChanged", accountChangedHandler)
        connector.removeListener("disconnect", disconnectHandler)
        connector.removeListener(
          "chooseWalletAndSendTransaction",
          showChooseWalletModal
        )
      }
    }
  }, [
    isConnected,
    connector,
    dispatch,
    yourAddress,
    web3,
    walletAddressFromUrl,
  ])
}

export default useSubscribeToConnectorEvents
