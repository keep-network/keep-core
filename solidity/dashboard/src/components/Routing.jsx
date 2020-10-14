import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import { withContractsDataContext } from "./ContractsDataContextProvider"
import Loadable from "./Loadable"
import { NotFound404 } from "./NotFound404"
import withWeb3Context from "./WithWeb3Context"
import OperatorPage from "../pages/OperatorPage"
import RewardsPageContainer from "../pages/RewardsPageContainer"
import CreateTokenGrantPage from "../pages/CreateTokenGrantPage"
import TokensPageContainer from "../pages/TokensPageContainer"
import ApplicationsPageContainer from "../pages/ApplicationsPageContainer"
import ResourcesPage from "../pages/ResourcesPage"
import TokenGrantPreviewPage from "../pages/TokenGrantPreviewPage"
import Header from "./Header"

const pages = [
  TokensPageContainer,
  OperatorPage,
  RewardsPageContainer,
  ApplicationsPageContainer,
]

const withoutWalletPages = [ResourcesPage]

class Routing extends React.Component {
  renderContent() {
    const {
      isKeepTokenContractDeployer,
      contractsDataIsFetching,
      web3: { error, provider, yourAddress },
    } = this.props

    if (!provider || !yourAddress) {
      // Temporary solution. We need to implement the empty states for other
      // pages.
      return <Header title="Choose wallet" />
    }

    if (error) {
      return null
    }

    return contractsDataIsFetching ? (
      <Loadable />
    ) : (
      <Switch>
        <Route exact path="/overview">
          {/* Temporary solution. We need to implement the `Overview` page as a separate page. */}
          <Redirect to="/tokens" />
        </Route>
        {pages.map(renderPage)}
        {isKeepTokenContractDeployer && (
          <Route
            exact
            path="/create-token-grants"
            component={CreateTokenGrantPage}
          />
        )}
        <Route exact path="/">
          <Redirect to="/tokens" />
        </Route>
        <Route path="*">
          <NotFound404 />
        </Route>
      </Switch>
    )
  }

  render() {
    return (
      <Switch>
        {withoutWalletPages.map(renderPage)}
        {/* In case that users will have bookmarked the old link. */}
        <Route exact path="/glossary">
          <Redirect to="/resources/quick-terminology" />
        </Route>
        <Route exact path="/grant/:grantId" component={TokenGrantPreviewPage} />
        <Route path="/">{this.renderContent()}</Route>
      </Switch>
    )
  }
}

export const renderPage = (PageComponent, index) => {
  return (
    <Route
      key={`${PageComponent.route.path}-${index}`}
      path={PageComponent.route.path}
      exact={PageComponent.route.exact}
    >
      <PageComponent
        routes={PageComponent.route.pages}
        title={PageComponent.route.title}
      />
    </Route>
  )
}

export default withWeb3Context(withContractsDataContext(Routing))
