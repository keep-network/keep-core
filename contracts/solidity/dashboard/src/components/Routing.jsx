import React from 'react'
import { BrowserRouter as Router, Route, Switch, Redirect } from 'react-router-dom'
import { Grid, Row, Col } from 'react-bootstrap'
import OverviewTab from './OverviewTab'
import StakeTab from './StakeTab'
import { RoutingTabs } from './RoutingTabs'
import TokenGrantsTab from './TokenGrantsTab'
import CreateTokenGrantsTab from './CreateTokenGrantsTab'
import { withContractsDataContext } from './ContractsDataContextProvider'
import Siginig from './Signing'
import Alerts from './Alerts'

class Routing extends React.Component { 

    renderRoutes = () => {
        const { isOperator, isTokenHolder } = this.props;

        if(isOperator)
            return <Route exact path='/overview' component={OverviewTab} />
        else if(!isOperator && !isTokenHolder)
            return <Redirect to='/sign-in' />
        return (
            <>
                <Route exact path='/overview' component={OverviewTab} />
                <Route exact path='/stake' component={StakeTab} />
                <Route exact path='/token-grants' component={TokenGrantsTab} />
                <Route exact path='/create-token-grants' component={CreateTokenGrantsTab} />
            </>
        )
    }

    renderContent() {
        const { isOperator, contractsDataIsFetching } = this.props
        if(contractsDataIsFetching)
            return <><div>Loading...</div></>
        
        return (
            <Switch>
                <Route exact path='/sign-in' component={Siginig} />
                <Route path="*">
                    <RoutingTabs isOperator={isOperator}>
                        <Switch>
                            <Redirect exact from='/' to='/overview' />
                            {this.renderRoutes()}
                        </Switch>
                    </RoutingTabs>
                </Route>
            </Switch>
            
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

export default withContractsDataContext(Routing)