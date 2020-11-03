import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import { NotFound404 } from "./NotFound404"
import { useWeb3Context } from "./WithWeb3Context"
import OperatorPage from "../pages/OperatorPage"
import RewardsPageContainer from "../pages/RewardsPageContainer"
import TokensPageContainer from "../pages/TokensPageContainer"
import ApplicationsPageContainer from "../pages/ApplicationsPageContainer"
import ResourcesPage from "../pages/ResourcesPage"
import TokenGrantPreviewPage from "../pages/TokenGrantPreviewPage"
import TokenOverviewPage from "../pages/OverviewPage"

const pages = [
  TokenOverviewPage,
  TokensPageContainer,
  OperatorPage,
  RewardsPageContainer,
  ApplicationsPageContainer,
  ResourcesPage,
]

class Routing extends React.Component {
  render() {
    return (
      <Switch>
        {/* In case that users will have bookmarked the old link. */}
        <Route exact path="/glossary">
          <Redirect to="/resources/quick-terminology" />
        </Route>
        <Route exact path="/grant/:grantId" component={TokenGrantPreviewPage} />
        {pages.map(renderPage)}
        <Route exact path="/">
          <Redirect to="/overview" />
        </Route>
        <Route path="*">
          <NotFound404 />
        </Route>
      </Switch>
    )
  }
}

export const renderPage = (PageComponent, index) => {
  return (
    <CustomRoute
      key={`${PageComponent.route.path}-${index}`}
      path={PageComponent.route.path}
      exact={PageComponent.route.exact}
      component={PageComponent}
      {...PageComponent.route}
    />
  )
}

const CustomRoute = ({
  path,
  exact,
  component: Component,
  emptyStateComponent: EmptyStateComponent = null,
  withConnectWalletGuard = false,
  ...componentProps
}) => {
  const { yourAddress, provider } = useWeb3Context()

  return !withConnectWalletGuard ||
    (withConnectWalletGuard && yourAddress && provider) ? (
    <Route path={path} exact={exact}>
      <Component routes={Component.route.pages} {...componentProps} />
    </Route>
  ) : (
    <EmptyStateComponent />
  )
}

export default Routing
