import LiquidityPage from "./LiquidityPage"
import EmptyStatePage from "./EmptyStatePage"

LiquidityPage.route = {
  title: "Liquidity",
  headerTitle: "Liquidity Rewards",
  path: "/liquidity",
  exact: true,
  withConnectWalletGuard: true,
  emptyStateComponent: EmptyStatePage,
  pages: [],
}

export default LiquidityPage
