import web3Utils from "web3-utils"
import { useRouteMatch } from "react-router-dom"
import useIsExactRoutePath from "./useIsExactRoutePath"

const useWalletAddressFromUrl = () => {
  const match = useRouteMatch(`/:address/:actualPath`)
  const isActualPathCorrect = useIsExactRoutePath(
    `/${match?.params?.actualPath}`,
    true
  )
  if (!match?.params?.address || !match?.params?.actualPath) return false

  if (isActualPathCorrect && web3Utils.isAddress(match.params.address)) {
    return match.params.address
  }

  return ""
}

export default useWalletAddressFromUrl
