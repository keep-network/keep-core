import web3Utils from "web3-utils"
import { useRouteMatch } from "react-router-dom"
import { pages } from "../components/Routing"

const useWalletAddressFromUrl = () => {
  const match = useRouteMatch(`/:address/:actualPath`)

  if (match) {
    const isActualPathCorrect = pages.some((page) => {
      return page.route.path === `/${match.params.actualPath}`
    })

    if (isActualPathCorrect && web3Utils.isAddress(match.params.address)) {
      return match.params.address
    }
  }

  return ""
}

export default useWalletAddressFromUrl
