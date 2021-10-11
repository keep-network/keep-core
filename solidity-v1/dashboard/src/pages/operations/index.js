import OperationsPage from "./OperatorPage"
import EmptyStatePage from "./EmptyStatePage"

OperationsPage.route = {
  title: "Operations",
  path: "/operations",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
  pages: [],
}

export default OperationsPage
