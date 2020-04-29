import React from "react"
import { Route, Switch, Redirect } from "react-router-dom"
import { withContractsDataContext } from "./ContractsDataContextProvider"
import Loadable from "./Loadable"
import { NotFound404 } from "./NotFound404"
import { withOnlyLoggedUser } from "./WithOnlyLoggedUserHoc"
import withWeb3Context from "./WithWeb3Context"
import TokensPage from "../pages/TokensPage"
import OperatorPage from "../pages/OperatorPage"
import AuthorizerPage from "../pages/AuthorizerPage"
import RewardsPage from "../pages/RewardsPage"
import CreateTokenGrantPage from "../pages/CreateTokenGrantPage"

class Routing extends React.Component {
  renderContent() {
    const {
      isKeepTokenContractDeployer,
      contractsDataIsFetching,
      web3: { error },
    } = this.props

    if (error) {
      return null
    }

    return contractsDataIsFetching ? (
      <Loadable />
    ) : (
      <Switch>
        <Route exact path="/tokens" component={TokensPage} />
        <Route exact path="/operations" component={OperatorPage} />
        <Route exact path="/rewards" component={RewardsPage} />
        <Route exact path="/authorizer" component={AuthorizerPage} />
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
    return <>{this.renderContent()}</>
  }
}

export default withWeb3Context(
  withContractsDataContext(withOnlyLoggedUser(Routing))
)
