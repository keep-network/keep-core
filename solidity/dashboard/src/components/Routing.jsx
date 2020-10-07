import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import { withContractsDataContext } from "./ContractsDataContextProvider"
import Loadable from "./Loadable"
import { NotFound404 } from "./NotFound404"
import withWeb3Context from "./WithWeb3Context"
import OperatorPage from "../pages/OperatorPage"
import RewardsPageContainer from "../pages/RewardsPageContainer"
import CreateTokenGrantPage from "../pages/CreateTokenGrantPage"
import TokenGrantsPage from "../pages/TokenGrantsPage"
import TokensPageContainer from "../pages/TokensPageContainer"
import ApplicationsPageContainer from "../pages/ApplicationsPageContainer"
import ChooseWallet from "./ChooseWallet"
import ResourcesPage from "../pages/ResourcesPage"
import TokenGrantPreviewPage from "../pages/TokenGrantPreviewPage"

class Routing extends React.Component {
  renderContent() {
    const {
      isKeepTokenContractDeployer,
      contractsDataIsFetching,
      web3: { error, provider, yourAddress },
    } = this.props

    if (!provider || !yourAddress) {
      return <ChooseWallet />
    }

    if (error) {
      return null
    }

    return contractsDataIsFetching ? (
      <Loadable />
    ) : (
      <Switch>
        <Route path="/tokens" component={TokensPageContainer} />
        <Route exact path="/operations" component={OperatorPage} />
        <Route path="/rewards" component={RewardsPageContainer} />
        <Route exact path="/token-grants" component={TokenGrantsPage} />
        <Route path="/applications" component={ApplicationsPageContainer} />
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
        <Route exact path="/resources" component={ResourcesPage} />
        {/* In case that users will have bookmarked the old link. */}
        <Route exact path="/glossary">
          <Redirect to="/resources#quick-terminology" />
        </Route>
        <Route exact path="/grant/:grantId" component={TokenGrantPreviewPage} />
        <Route path="/">{this.renderContent()}</Route>
      </Switch>
    )
  }
}

export default withWeb3Context(withContractsDataContext(Routing))
