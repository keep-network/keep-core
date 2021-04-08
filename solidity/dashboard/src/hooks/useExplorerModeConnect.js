import React, { useEffect } from "react"
import web3Utils from "web3-utils"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import ExplorerModeModal from "../components/ExplorerModeModal"
import { useLocation } from "react-router-dom"
import { useModal } from "./useModal"
import { useWeb3Context } from "../components/WithWeb3Context"

/**
 * Checks if there is addres in the url and the tries to connect to explorer mode
 *
 * Url pattern: http(s)://<site_name>/<address>/<page>
 * Example for localhost:
 *  http://localhost:3000/0x857173e7c7d76e051e80d30FCc3EA6A9C2b53756/liquidity
 *
 */
const useExplorerModeConnect = () => {
  const location = useLocation()
  const { openModal, closeModal } = useModal()
  const { connector, connectAppWithWallet } = useWeb3Context()

  useEffect(() => {
    const pathnameSplitted = location.pathname.split("/")
    if (pathnameSplitted.length > 1 && pathnameSplitted[1] && !connector) {
      const address = pathnameSplitted[1]
      if (web3Utils.isAddress(address)) {
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
            title: "Connect Wallet",
          }
        )
      }
    }
  }, [location.pathname])
}

export default useExplorerModeConnect
