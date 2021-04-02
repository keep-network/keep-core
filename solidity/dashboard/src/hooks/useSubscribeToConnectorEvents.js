import React, { useEffect } from "react"
import { useDispatch } from "react-redux"
import { useWeb3Context } from "../components/WithWeb3Context"
import { WALLETS } from "../constants/constants"
import { useModal } from "./useModal"
import { WalletSelectionModal } from "../components/WalletSelectionModal"

const useSubscribeToConnectorEvents = () => {
  const dispatch = useDispatch()
  const { isConnected, connector, yourAddress } = useWeb3Context()
  const { openModal, closeModal } = useModal()

  useEffect(() => {
    const accountChangedHandler = (address) => {
      dispatch({ type: "app/account_changed", payload: { address } })
    }

    const disconnectHandler = () => {
      dispatch({ type: "app/logout" })
    }

    const showChooseWalletModal = (payload) => {
      openModal(<WalletSelectionModal />)
    }

    if (isConnected && connector) {
      dispatch({ type: "app/login", payload: { address: yourAddress } })
      if (connector.name === WALLETS.METAMASK.name) {
        connector.getProvider().on("accountsChanged", accountChangedHandler)
        connector.getProvider().on("chainChanged", disconnectHandler)
      }

      if (connector.name === WALLETS.READ_ONLY_ADDRESS.name) {
        connector.eventEmitter.on(
          "chooseWalletAndSendTransaction",
          showChooseWalletModal
        )
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
  }, [isConnected, connector, dispatch, yourAddress])
}

export default useSubscribeToConnectorEvents
