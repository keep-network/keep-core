import Web3 from 'web3'
import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'

import './App.css';

import Token from './components/Token';
import Staking from './components/Staking';
import Vesting from './components/Vesting';

const App = () => (
  <Router>
    <Switch>
      <Route path="/token/:address" component={ Main }/>
      <Route path="/staking/:address" component={ Staking }/>
      <Route path="/vesting/:address" component={ Vesting }/>
      <Route component={ MissingAddress } />
    </Switch>
  </Router>
)

const Main = function({ match }) {
  let web3 = new Web3()
  let { address } = match.params

  return web3.utils.isAddress(address)
    ? <Token address={ address } />
    : <MissingAddress />
}

const MissingAddress = () => (
  <p>This is not a Token address</p>
)

export default App;
