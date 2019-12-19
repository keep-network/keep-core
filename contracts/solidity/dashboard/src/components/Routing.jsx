import React from 'react'
import { BrowserRouter as Router, Route, Switch, Redirect } from 'react-router-dom'
import { Grid, Row, Col } from 'react-bootstrap'
import OverviewTab from './OverviewTab'
import StakeTab from './StakeTab'
import { RoutingTabs } from './RoutingTabs'
import TokenGrantsTab from './TokenGrantsTab'
import CreateTokenGrantsTab from './CreateTokenGrantsTab'
import { withContractsDataContext } from './ContractsDataContextProvider'
import Alerts from './Alerts'
import Loadable from './Loadable'
import { NotFound404 } from './NotFound404'
import { WithOnlyLoggedUser } from './WithOnlyLoggedUserHoc'
import WithWeb3Context from './WithWeb3Context'

class Routing extends React.Component { 

    renderContent() {
        const { isOperator, isTokenHolder, contractsDataIsFetching, web3: { error }} = this.props


        if(error) {
            return null;
        }
        return contractsDataIsFetching ? <Loadable /> : (
            <RoutingTabs isOperator={isOperator} isTokenHolder={isTokenHolder}>
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
            </RoutingTabs>
        )
    }

    render() {
        return (
           <Router>
               <Grid>
                   <Row>
                        <Col xs={12}>
                            <Alerts />
                            {this.renderContent()}
                        </Col>
                   </Row>
               </Grid>
           </Router>
       );
    }
}

export default WithWeb3Context(withContractsDataContext(WithOnlyLoggedUser(Routing)))