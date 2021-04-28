import React, { useEffect, useState } from "react"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import ExplorerModeModal from "../components/ExplorerModeModal"
import { useLocation, useHistory } from "react-router-dom"
import { useModal } from "./useModal"
import { useWeb3Context } from "../components/WithWeb3Context"
import { WALLETS } from "../constants/constants"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl"
import { injected } from "../connectors"
import { isEmptyArray } from "../utils/array.utils"
import { usePrevious } from "./usePrevious"

const useHasChanged = (val) => {
  const prevVal = usePrevious(val)
  return prevVal !== val
}

/**
 * Checks if there is a wallet addres in the url and then tries to connect to
 * Explorer Mode.
 *
 * Also changes the url after connecting and disconnecting from the Eplorer Mode
 *
 * Url pattern: http(s)://<site_name>/<address>/<page>
 * Example for localhost:
 *  http://localhost:3000/0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756/liquidity
 *
 */
const useAutoConnect = () => {
  const location = useLocation()
  const history = useHistory()
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const { openModal, closeModal } = useModal()
  const { connector, connectAppWithWallet, yourAddress } = useWeb3Context()

  const locationHasChanged = useHasChanged(location.pathname)
  const [injectedTried, setInjectedTried] = useState(false)

  useEffect(() => {
    if (locationHasChanged) return
    // change url to the one with address when we connect to the explorer mode
    if (!walletAddressFromUrl && connector && yourAddress) {
      const newPathname = "/" + yourAddress + location.pathname
      history.push({ pathname: newPathname })
    }
  }, [
    connector,
    yourAddress,
    history,
    location.pathname,
    locationHasChanged,
    walletAddressFromUrl,
  ])

  useEffect(() => {
    if (injectedTried) return
    injected.getAccounts().then((accounts) => {
      setInjectedTried(true)
      if (!isEmptyArray(accounts) && !walletAddressFromUrl) {
        connectAppWithWallet(injected, false).catch((error) => {
          // Just log an error, we don't want to do anything else.
          console.log(
            "Eager injected connector cannot connect with the dapp:",
            error.message
          )
        })
      } else if (walletAddressFromUrl && !connector) {
        const explorerModeConnector = new ExplorerModeConnector()
        openModal(
          <ExplorerModeModal
            connectAppWithWallet={connectAppWithWallet}
            connector={explorerModeConnector}
            closeModal={closeModal}
            address={walletAddressFromUrl}
            connectWithWalletOnMount={true}
          />,
          {
            title: "Connect Ethereum Address",
          }
        )
      }
    })
  }, [
    connectAppWithWallet,
    walletAddressFromUrl,
    injectedTried,
    closeModal,
    openModal,
    connector,
  ])
}

export default useAutoConnect
