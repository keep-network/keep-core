import React from 'react'
import { Route, Switch, Redirect } from 'react-router-dom'
import OverviewTab from './OverviewTab'
import StakeTab from './StakeTab'
import TokenGrantsTab from './TokenGrantsTab'
import CreateTokenGrantsTab from './CreateTokenGrantsTab'
import { withContractsDataContext } from './ContractsDataContextProvider'
import Alerts from './Alerts'
import Loadable from './Loadable'
import { NotFound404 } from './NotFound404'
import { withOnlyLoggedUser } from './WithOnlyLoggedUserHoc'
import withWeb3Context from './WithWeb3Context'
import { Rewards } from './Rewards'
import TokensPage from '../pages/TokensPage'

class Routing extends React.Component {
  renderContent() {
    const { isTokenHolder, isKeepTokenContractDeployer, contractsDataIsFetching, web3: { error } } = this.props

    if (error) {
      return null
    }

    return contractsDataIsFetching ? <Loadable /> : (
      <Switch>
        <Route exact path='/tokens' component={TokensPage} />
        <Route exact path='/rewards' component={Rewards} />
        {isTokenHolder && <Route exact path='/stake' component={StakeTab} /> }
        {isTokenHolder && <Route exact path='/token-grants' component={TokenGrantsTab} /> }
        {isKeepTokenContractDeployer && <Route exact path='/create-token-grants' component={CreateTokenGrantsTab} /> }
        <Route exact path='/' >
          <Redirect to='/overview' />
        </Route>
        <Route path="*">
          <NotFound404 />
        </Route>
      </Switch>
    )
  }

  render() {
    return (
      <>
        <Alerts />
        {this.renderContent()}
      </>
    )
  }
}

export default withWeb3Context(withContractsDataContext(withOnlyLoggedUser(Routing)))
