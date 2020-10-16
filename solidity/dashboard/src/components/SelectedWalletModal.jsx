import React, { useEffect, useState } from "react"
import { KeepLoadingIndicator } from "./Icons"
import { useShowMessage, messageType } from "./Message"
import ChooseWalletAddress from "./ChooseWalletAddress"
import { isEmptyArray } from "../utils/array.utils"
import { wait } from "../utils/general.utils"

const SelectedWalletModal = ({
  icon,
  walletName,
  iconDescription,
  description,
  providerName,
  connector,
  connectAppWithWallet,
  closeModal,
  fetchAvailableAccounts = null,
  numberOfAccounts = 15,
  connectWithWalletOnMount = false,
  withAccountPagination = false,
  children,
}) => {
  const showMessage = useShowMessage()
  const [isConnecting, setIsConnecting] = useState(false)
  const [accountsOffSet, setAccountsOffSet] = useState(0)
  const [accountsAreFetching, setAccountsAreFetching] = useState(false)
  const [availableAccounts, setAvailableAccounts] = useState([])
  const [error, setError] = useState("")

  useEffect(() => {
    let shouldSetState = true
    // Fetching wallet addresses.
    if (
      fetchAvailableAccounts &&
      typeof fetchAvailableAccounts === "function" &&
      connector
    ) {
      setAccountsAreFetching(true)
      fetchAvailableAccounts(numberOfAccounts, accountsOffSet)
        .then((availableAccounts) => {
          if (shouldSetState) {
            setAvailableAccounts(availableAccounts)
            setAccountsAreFetching(false)
          }
        })
        .catch((error) => {
          if (shouldSetState) {
            handleError(error)
            setAccountsAreFetching(false)
          }
        })
    }
    return () => {
      shouldSetState = false
    }
  }, [connector, fetchAvailableAccounts, numberOfAccounts, accountsOffSet])

  useEffect(() => {
    let shouldSetState = true
    // Connecting to a wallet when a component did mount.
    if (connector && connectWithWalletOnMount) {
      setIsConnecting(true)
      wait(1000) // Delay request and show loading indicator.
        .then(() => connectAppWithWallet(connector, providerName))
        .then(() => {
          if (shouldSetState) {
            setIsConnecting(false)
          }
          closeModal()
        })
        .catch((error) => {
          if (shouldSetState) {
            handleError(error)
          }
          showMessage({
            type: messageType.ERROR,
            title: error.message,
            sticky: true,
          })
        })
    }

    return () => {
      shouldSetState = false
    }
  }, [
    connector,
    connectAppWithWallet,
    providerName,
    closeModal,
    showMessage,
    connectWithWalletOnMount,
  ])

  const handleError = (error) => {
    setError(error.toString())
    setIsConnecting(false)
  }

  const onSelectAccount = async (account) => {
    try {
      connector.defaultAccount = account
      setIsConnecting(true)
      await connectAppWithWallet(connector, providerName)
      setIsConnecting(false)
      closeModal()
    } catch (error) {
      handleError(error)
      showMessage({ type: messageType.ERROR, title: error.message })
    }
  }

  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        {icon}
        <h3 className="ml-1">{walletName}</h3>
      </div>
      {iconDescription && iconDescription}
      <span className="text-center">{description}</span>
      {children}
      {isConnecting || accountsAreFetching ? (
        <>
          <KeepLoadingIndicator />
          {accountsAreFetching
            ? `Loading wallet addresses...`
            : `Connecting...`}
        </>
      ) : null}
      {error && error}
      {!isEmptyArray(availableAccounts) &&
        !accountsAreFetching &&
        !isConnecting && (
          <ChooseWalletAddress
            onSelectAccount={onSelectAccount}
            addresses={availableAccounts}
            withPagination={withAccountPagination}
            renderPrevBtn={accountsOffSet > 0}
            onNext={() => setAccountsOffSet((prevOffset) => prevOffset + 5)}
            onPrev={() => setAccountsOffSet((prevOffset) => prevOffset - 5)}
          />
        )}
    </div>
  )
}

export default React.memo(SelectedWalletModal)
