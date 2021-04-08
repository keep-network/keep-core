import { useLocation } from "react-router-dom"
import web3Utils from "web3-utils"

const useWalletAddressFromUrl = () => {
  const location = useLocation()

  const pathnameSplitted = location.pathname.split("/")

  if (pathnameSplitted.length > 1 && pathnameSplitted[1]) {
    const address = pathnameSplitted[1]
    if (web3Utils.isAddress(address)) {
      return address
    }
  }

  return ""
}

export default useWalletAddressFromUrl
