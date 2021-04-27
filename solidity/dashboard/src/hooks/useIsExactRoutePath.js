import { pages } from "../components/Routing"

const useIsExactRoutePath = (path) => {
  if (!path) return false

  return pages.some((page) => {
    let subpageMatch = false
    // check if subpages match the path
    if (page.route.pages?.length > 0) {
      subpageMatch = page.route.pages.some((subpage) => {
        return subpage.route.path === path
      })
    }
    return subpageMatch ? subpageMatch : page.route.path === path
  })
}

export default useIsExactRoutePath
