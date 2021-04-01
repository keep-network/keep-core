import { useEffect } from "react"
import { useWeb3Context } from "../components/WithWeb3Context"
import { injected } from "../connectors"
import { isEmptyArray } from "../utils/array.utils"

const useInjectedAutoConnect = () => {
  const { connectAppWithWallet } = useWeb3Context()

  useEffect(() => {
    injected.getAccounts().then((accounts) => {
      if (!isEmptyArray(accounts)) {
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

export default useInjectedAutoConnect
