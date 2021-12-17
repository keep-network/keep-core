import { useEffect, useState, useCallback } from "react"
import { ExplorerModeConnector } from "../connectors/explorer-mode-connector"
import { useModal } from "./useModal"
import { useWeb3Context } from "../components/WithWeb3Context"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl"
import {
  metaMaskInjectedConnector,
  tallyInjectedConnector,
} from "../connectors"
import { isEmptyArray } from "../utils/array.utils"
import useIsExactRoutePath from "./useIsExactRoutePath"
import { isSameEthAddress } from "../utils/general.utils"
import { MODAL_TYPES } from "../constants/constants"

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
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const { openModal } = useModal()
  const { connector, connectAppWithWallet } = useWeb3Context()
  const isExactRoutePath = useIsExactRoutePath()
  const [injectedTried, setInjectedTried] = useState(false)

  const isWalletFromUrlSameAsInProvider = useCallback(
    (providerAccounts) => {
      return (
        walletAddressFromUrl &&
        providerAccounts.some((account) =>
          isSameEthAddress(account, walletAddressFromUrl)
        )
      )
    },
    [walletAddressFromUrl]
  )

  useEffect(() => {
    if (injectedTried) return

    const connectWallet = (connector) => {
      connector.getAccounts().then((accounts) => {
        setInjectedTried(true)
        if (
          (!isEmptyArray(accounts) && isExactRoutePath) ||
          isWalletFromUrlSameAsInProvider(accounts)
        ) {
          connectAppWithWallet(connector, false).catch((error) => {
            // Just log an error, we don't want to do anything else.
            console.log(
              "Eager injected connector cannot connect with the dapp:",
              error.message
            )
          })
        } else if (walletAddressFromUrl && !connector) {
          const explorerModeConnector = new ExplorerModeConnector()
          openModal(MODAL_TYPES.ExplorerMode, {
            connectAppWithWallet,
            connector: explorerModeConnector,
            address: walletAddressFromUrl,
            connectWithWalletOnMount: true,
          })
        }
      })
    }

    if (window.tally) {
      connectWallet(tallyInjectedConnector)
    } else if (window.ethereum && window.ethereum.isMetaMask) {
      connectWallet(metaMaskInjectedConnector)
    }
  }, [
    connectAppWithWallet,
    walletAddressFromUrl,
    injectedTried,
    openModal,
    connector,
    isExactRoutePath,
    isWalletFromUrlSameAsInProvider,
  ])
}

export default useAutoConnect
