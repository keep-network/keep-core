import React from 'react'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import Main from './components/Main';
import Web3ContextProvider from './components/Web3ContextProvider';

const App = () => (
  <Router>
    <Switch>
      <Web3ContextProvider>
        <Route component={ Main } />
      </Web3ContextProvider>
    </Switch>
  </Router>
)

export default App