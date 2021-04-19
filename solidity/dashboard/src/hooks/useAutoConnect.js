import React, { useEffect } from "react"
import web3Utils from "web3-utils"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import ExplorerModeModal from "../components/ExplorerModeModal"
import { useLocation, useHistory } from "react-router-dom"
import { useModal } from "./useModal"
import { useWeb3Context } from "../components/WithWeb3Context"
import { WALLETS } from "../constants/constants"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl"
import { injected } from "../connectors"
import { isEmptyArray } from "../utils/array.utils"

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

  useEffect(() => {
    const pathnameSplitted = location.pathname.split("/")
    if (pathnameSplitted.length > 1 && pathnameSplitted[1]) {
      // change url to the one without an address when disconnecting
      if (web3Utils.isAddress(walletAddressFromUrl) && !connector) {
        const newPathname = location.pathname.replace(
          "/" + walletAddressFromUrl,
          ""
        )
        history.push({ pathname: newPathname })
      }

      // change url to the one with address when we connect to the explorer mode
      if (
        !web3Utils.isAddress(walletAddressFromUrl) &&
        connector &&
        yourAddress &&
        connector.name === WALLETS.EXPLORER_MODE.name
      ) {
        const newPathname = "/" + yourAddress + location.pathname
        history.push({ pathname: newPathname })
      }
    }
  }, [connector, yourAddress])

  useEffect(() => {
    // history.push({ pathname: "elo/melo" })
    const pathnameSplitted = location.pathname.split("/")
    if (pathnameSplitted.length > 1 && pathnameSplitted[1]) {
      // log in to explorer mode when pasting url with address
      if (web3Utils.isAddress(walletAddressFromUrl) && !connector) {
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
    }
  }, [location.pathname])

  useEffect(() => {
    injected.getAccounts().then((accounts) => {
      if (!isEmptyArray(accounts) && !walletAddressFromUrl) {
        connectAppWithWallet(injected, false).catch((error) => {
          // Just log an error, we don't want to do anything else.
          console.log(
            "Eager injected connector cannot connect with the dapp:",
            error.message
          )
        })
      }
    })
  }, [connectAppWithWallet])
}

export default useAutoConnect
