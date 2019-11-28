import React from 'react';
import { BrowserRouter as Router, Route, Switch, Redirect, useHistory } from 'react-router-dom'
import { Grid, Row, Col } from 'react-bootstrap';
import OverviewTab from './OverviewTab';
import StakeTab from './StakeTab';
import { RoutingTabs } from './RoutingTabs'
import TokenGrantsTab from './TokenGrantsTab';
import CreateTokenGrantsTab from './CreateTokenGrantsTab';
import { withContractsDataContext } from './ContractsDataContextProvider';
import SigningForm from './SigningForm'
import Siginig from './Signing';

class Routing extends React.Component { 

    renderContent = () => {
        const { isOperator, isTokenHolder, contractsDataIsFetching } = this.props;

        if(contractsDataIsFetching)
            return <><div>Loading...</div></>
        else if(isOperator)
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

    render() {
        console.log('props', this.props);
        return (
           <Router>
               <Grid>
                   <Row>
                        <Col xs={12}>
                        <RoutingTabs>
                            <Switch>
                                <Route exact path='/sign-in' component={Siginig} />
                                {this.renderContent()}
                                <Redirect to='/overview'/>
                            </Switch>
                        </RoutingTabs>
                        </Col>
                   </Row>
               </Grid>
           </Router>
       );
    }
}

export default withContractsDataContext(Routing)