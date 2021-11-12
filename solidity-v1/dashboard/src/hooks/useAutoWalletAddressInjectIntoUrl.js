import { useEffect } from "react"
import { useHistory, useLocation } from "react-router-dom"
import useWalletAddressFromUrl from "./useWalletAddressFromUrl"
import { useWeb3Context } from "../components/WithWeb3Context"
import useHasChanged from "./useHasChanged"
import { isSameEthAddress } from "../utils/general.utils"

const useAutoWalletAddressInjectIntoUrl = () => {
  const location = useLocation()
  const history = useHistory()
  const walletAddressFromUrl = useWalletAddressFromUrl()
  const { connector, yourAddress } = useWeb3Context()

  const locationHasChanged = useHasChanged(location.pathname)

  useEffect(() => {
    if (locationHasChanged) return
    // change url to the one with address when we connect to the explorer mode
    if (!walletAddressFromUrl && connector && yourAddress) {
      const newPathname = "/" + yourAddress + location.pathname
      history.push({ pathname: newPathname })
    }

    if (!walletAddressFromUrl && connector && yourAddress) {
      const newPathname = "/" + yourAddress + location.pathname
      history.push({ pathname: newPathname })
    } else if (walletAddressFromUrl && connector && yourAddress) {
      if (!isSameEthAddress(yourAddress, walletAddressFromUrl)) {
        const newPathname =
          "/" +
          yourAddress +
          location.pathname.replace(`/${walletAddressFromUrl}`, "")
        history.push({ pathname: newPathname })
      }
    }
  }, [
    connector,
    yourAddress,
    history,
    location.pathname,
    locationHasChanged,
    walletAddressFromUrl,
  ])
}

export default useAutoWalletAddressInjectIntoUrl
