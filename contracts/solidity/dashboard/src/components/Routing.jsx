import React from 'react';
import { BrowserRouter as Router, Route, Switch, Redirect, useHistory } from 'react-router-dom'
import { Grid, Row, Col } from 'react-bootstrap';
import OverviewTab from './OverviewTab';
import StakeTab from './StakeTab';
import { RoutingTabs } from './RoutingTabs'
import TokenGrantsTab from './TokenGrantsTab';
import CreateTokenGrantsTab from './CreateTokenGrantsTab';

export default class Routing extends React.Component { 

    render() {
        return (
           <Router>
               <Grid>
                   <Row>
                        <Col xs={12}>
                        <RoutingTabs>
                            <Switch>
                                <Route exact path='/overview' component={OverviewTab} />
                                <Route exact path='/stake' component={StakeTab} />
                                <Route exact path='/token-grants' component={TokenGrantsTab} />
                                <Route exact path='/create-token-grants' component={CreateTokenGrantsTab} />
                                <Redirect from='/' to='/overview'/>
                            </Switch>
                        </RoutingTabs>
                        </Col>
                   </Row>
               </Grid>
           </Router>
       );
    }
}