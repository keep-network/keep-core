import { EmptyStatePage } from "./EmptyStatePage"
import TokenGrantsPage, { TokenGrantPreviewPage } from "./TokenGrantsPage"

TokenGrantsPage.route = {
  title: "Token Grants",
  path: "/token-grants",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
}

TokenGrantPreviewPage.route = {
  title: "Token Grants",
  path: "/grant/:grantId",
  exact: true,
}

export default TokenGrantsPage
export { TokenGrantPreviewPage }
