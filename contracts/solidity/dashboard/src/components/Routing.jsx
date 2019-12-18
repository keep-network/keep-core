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

class Routing extends React.Component { 

    renderContent() {
        const { isOperator, isTokenHolder, contractsDataIsFetching } = this.props

        return contractsDataIsFetching ? <Loadable /> : (
            <Switch>
                <Route exact path='/overview' component={OverviewTab} />
                {isTokenHolder && <Route exact path='/stake' component={StakeTab} /> }
                {isTokenHolder && <Route exact path='/token-grants' component={TokenGrantsTab} /> }
                {isTokenHolder && <Route exact path='/create-token-grants' component={CreateTokenGrantsTab} /> }
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
           
       );
    }
}

export default withContractsDataContext(Routing)