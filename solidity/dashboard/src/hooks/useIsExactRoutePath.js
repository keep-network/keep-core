import { pages } from "../components/Routing"
import { useLocation } from "react-router-dom"
import { removeSubstringBetweenCharacter } from "../utils/general.utils"

const useIsExactRoutePath = (path = null, cutOutWalletAddress = false) => {
  const location = useLocation()
  const realPath = path
    ? path
    : cutOutWalletAddress
    ? removeSubstringBetweenCharacter(location.pathname, "/", 0)
    : location.pathname

  if (!realPath) return false

  return pages.some((page) => {
    let subpageMatch = false
    // check if subpages match the path
    if (page.route.pages?.length > 0) {
      subpageMatch = page.route.pages.some((subpage) => {
        return subpage.route.path === realPath
      })
    }
    return subpageMatch ? subpageMatch : page.route.path === realPath
  })
}

export default useIsExactRoutePath
