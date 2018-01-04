import React, { Component } from 'react';
import { Table, Col, Grid, Row } from 'react-bootstrap';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { displayAmount } from './utils';
import Network from './network';
import { getKeepToken, getTokenStaking } from './contracts';
import Header from './components/Header';
import BalanceChart from './components/BalanceChart';
import './App.css';

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
    this.state = {}
  }

  componentDidMount() {
    this.getData();
  }

  render() {
    const { yourAddress, tokenBalance, stakeBalance } = this.state;
    return (
      <div className="main">
        <Header />
        <Grid>
          <Row>
            <Col xs={12} md={12}>
              <BalanceChart/>
              <Table className="small" striped bordered condensed>
                <tbody>
                  <TableRow title="Your address">
                    { yourAddress }
                  </TableRow>
                  <TableRow title="Token Balance">
                    { tokenBalance } 
                  </TableRow>
                  <TableRow title="Stake Balance">
                    { stakeBalance } 
                  </TableRow>
                </tbody>
              </Table>
            </Col>  
            <Col xs={12} md={12}>
            </Col>
          </Row>
        </Grid>
      </div>
    )
  }

  async getData() {
    const accounts = await Network.getAccounts();
    const yourAddress  = accounts[0];
    const token = await getKeepToken(process.env.REACT_APP_TOKEN_ADDRESS);
    const tokenBalance =  displayAmount(await token.balanceOf(yourAddress), 3);
    
    const stakingContract = await getTokenStaking(process.env.REACT_APP_STAKING_ADDRESS);
    const stakeBalance  = displayAmount(await stakingContract.balanceOf(yourAddress), 3);
    
    this.setState({
      yourAddress,
      tokenBalance,
      stakeBalance
    })
  }
}

function TableRow({ title, children }) {
  return (
    <tr>
      <th><strong>{ title }</strong></th>
      <td>
        { children }
      </td>
    </tr>
  )
}

export default App;
