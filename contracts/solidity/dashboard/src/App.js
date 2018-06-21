import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import BigNumber from "bignumber.js";
import { displayAmount } from './utils';
import Network from './network';
import { getKeepToken, getTokenStaking, getTokenGrant } from './contracts';

const App = () => (
  <Router>
    <Switch>
      <Route component={ Main } />
    </Switch>
  </Router>
)

class Main extends Component {

  constructor() {
    super()
    this.state = {};
  }

  render() {
    return (
      <div className="main">
      </div>
    )
  }
}
export default App;
