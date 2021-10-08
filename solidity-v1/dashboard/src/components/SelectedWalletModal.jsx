import React, { useEffect, useState } from "react"
import { KeepLoadingIndicator } from "./Icons"
import { useShowMessage, messageType } from "./Message"
import ChooseWalletAddress from "./ChooseWalletAddress"
import { isEmptyArray } from "../utils/array.utils"
import { wait } from "../utils/general.utils"
import { UserRejectedConnectionRequestError } from "../connectors"

const SelectedWalletModal = ({
  icon,
  walletName,
  descriptionIcon,
  description,
  connector,
  connectAppWithWallet,
  closeModal,
  userRejectedConnectionRequestErrorMsg,
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
        .then(() => connectAppWithWallet(connector))
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
            messageType: messageType.ERROR,
            messageProps: {
              content: error.message,
              sticky: true,
            },
          })
        })
    }

    return () => {
      shouldSetState = false
    }
  }, [
    connector,
    connectAppWithWallet,
    closeModal,
    showMessage,
    connectWithWalletOnMount,
  ])

  const handleError = (error) => {
    console.error("Failed to connect to a wallet.", error)
    setError(error)
    setIsConnecting(false)
  }

  const onSelectAccount = async (account) => {
    try {
      connector.defaultAccount = account
      setIsConnecting(true)
      await connectAppWithWallet(connector)
      setIsConnecting(false)
      closeModal()
    } catch (error) {
      handleError(error)
      showMessage({
        messageType: messageType.ERROR,
        messageProps: {
          content: error.message,
          sticky: true,
        },
      })
    }
  }

  const renderError = () => {
    const parseError = (msg) => `Error: ${msg}`
    if (!error) {
      return null
    }

    if (error && error instanceof UserRejectedConnectionRequestError) {
      return parseError(userRejectedConnectionRequestErrorMsg || error.message)
    }

    return error && error.message
      ? parseError(error.message.toString())
      : parseError("Unexpected error, please try again.")
  }

  return (
    <div className="flex column center">
      <div className="flex full-center mb-3">
        {icon}
        <h3 className="ml-1">{walletName}</h3>
      </div>
      {descriptionIcon && descriptionIcon}
      <span className="text-center mt-1">{description}</span>
      {children}
      {isConnecting || accountsAreFetching ? (
        <>
          <KeepLoadingIndicator />
          {accountsAreFetching
            ? `Loading wallet addresses...`
            : `Connecting...`}
        </>
      ) : null}
      {renderError()}
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
