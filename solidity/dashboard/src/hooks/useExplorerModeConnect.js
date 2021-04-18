import React, { useEffect } from "react"
import web3Utils from "web3-utils"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import ExplorerModeModal from "../components/ExplorerModeModal"
import { useLocation, useHistory } from "react-router-dom"
import { useModal } from "./useModal"
import { useWeb3Context } from "../components/WithWeb3Context"
import { WALLETS } from "../constants/constants"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl";

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
const useExplorerModeConnect = () => {
  const location = useLocation()
  const history = useHistory()
  const address = useWalletAddressFromUrl()
  const { openModal, closeModal } = useModal()
  const { connector, connectAppWithWallet, yourAddress } = useWeb3Context()

  useEffect(() => {
    const pathnameSplitted = location.pathname.split("/")
    if (pathnameSplitted.length > 1 && pathnameSplitted[1]) {
      // change url to the one without an address when disconnecting
      if (web3Utils.isAddress(address) && !connector) {
        const newPathname = location.pathname.replace("/" + address, "")
        history.push({ pathname: newPathname })
      }

      // change url to the one with address when we connect to the explorer mode
      if (
        !web3Utils.isAddress(address) &&
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
    const pathnameSplitted = location.pathname.split("/")
    if (pathnameSplitted.length > 1 && pathnameSplitted[1]) {
      // log in to explorer mode when pasting url with address
      if (web3Utils.isAddress(address) && !connector) {
        const explorerModeConnector = new ExplorerModeConnector()
        openModal(
          <ExplorerModeModal
            connectAppWithWallet={connectAppWithWallet}
            connector={explorerModeConnector}
            closeModal={closeModal}
            address={address}
            connectWithWalletOnMount={true}
          />,
          {
            title: "Connect Ethereum Address",
          }
        )
      }
    }
  }, [location.pathname])
}

export default useExplorerModeConnect
