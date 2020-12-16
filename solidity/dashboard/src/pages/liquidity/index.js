import LiquidityPage from "./LiquidityPage"
import EmptyStatePage from "./EmptyStatePage"

LiquidityPage.route = {
  title: "Liquidity",
  path: "/liquidity",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
  pages: [],
}

export default LiquidityPage
