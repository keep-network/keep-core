import React, { Component } from 'react';
import { Pie } from 'react-chartjs-2';
import { Table, Col, Grid, Row } from 'react-bootstrap';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import moment from 'moment';
import BigNumber from "bignumber.js";
import { displayAmount } from './utils';
import Network from './network';
import { getKeepToken, getTokenStaking, getTokenGrant } from './contracts';
import Header from './components/Header';
import StakingForm from './components/StakingForm';
import WithdrawalsTable from './components/WithdrawalsTable';

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
    this.state.chartData = {};
  }

  componentDidMount() {
    this.getData();
  }

  render() {
    const { yourAddress, tokenBalance, stakeBalance, grantBalance, grantStakeBalance, chartOptions, chartData, withdrawals, withdrawalsTotal } = this.state;

    return (
      <div className="main">
        <Header />
        <Grid>
          <Row>
            <Col xs={12} md={6}>
              <Pie dataKey="name" data={ chartData } options={ chartOptions } />
            </Col>
            <Col xs={12} md={6}>
              
              <Table className="small table-sm" striped bordered condensed>
                <tbody>
                  <TableRow title="Your wallet address">
                    { yourAddress }
                  </TableRow>
                  <TableRow title="Tokens">
                    { tokenBalance } 
                  </TableRow>
                  <TableRow title="Staked">
                    { stakeBalance } 
                  </TableRow>
                  <TableRow title="Pending withdrawals">
                    { withdrawalsTotal } 
                  </TableRow>
                  <TableRow title="Token Grants">
                    { grantBalance } 
                  </TableRow>
                  <TableRow title="Staked Token Grants">
                    { grantStakeBalance } 
                  </TableRow>
                </tbody>
              </Table>

              <StakingForm btnText="Stake" action="stake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>
              <StakingForm btnText="Unstake" action="unstake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>

            </Col>  
            
          </Row>
          <Row>
          <Col xs={12} md={4}>
              <h4>Pending withdrawals</h4>
              <WithdrawalsTable data={withdrawals}/>
            </Col>  
            <Col xs={12} md={4}>
              <h4>Token grants</h4>
              <Table className="small table-sm" condensed>
                <thead>
                  <tr>
                    <th><strong>Amount</strong></th>
                  </tr>
                </thead>
                <tbody>
                </tbody>
              </Table>
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

    const tokenBalance =  displayAmount(await token.balanceOf(yourAddress), 18, 3);

    const stakingContract = await getTokenStaking(process.env.REACT_APP_STAKING_ADDRESS);
    const stakeBalance  = displayAmount(await stakingContract.balanceOf(yourAddress), 18, 3);
    const withdrawalDelay  = (await stakingContract.withdrawalDelay()).toNumber()

    const withdrawalIndexes  = await stakingContract.getWithdrawals(yourAddress);
    let withdrawals = [];
    let withdrawalsTotal = 0;
    // const now = new Date() / 1000; // normalize to seconds

    for(let i=0; i < withdrawalIndexes.length; i++) {
      const withdrawal  = await stakingContract.getWithdrawal(withdrawalIndexes[i].toNumber());
      const availableAt = moment((withdrawalDelay+withdrawal[2].toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a");
      withdrawals.push({
        'amount': withdrawal[1].toNumber(),
        'availableAt': availableAt
        }
      );
      withdrawalsTotal += withdrawal[1].toNumber();
    }

    const grantContract = await getTokenGrant(process.env.REACT_APP_TOKENGRANT_ADDRESS);
    const grantBalance  = displayAmount(await grantContract.balanceOf(yourAddress), 18, 3);
    const grantStakeBalance  = displayAmount(await grantContract.stakeBalanceOf(yourAddress), 18, 3);

    const chartOptions = {
      legend: {
        position: 'right'
      }
    }
    const chartData = {
      labels: [
        'Tokens',
        'Staked',
        'Pending withdrawals',
        'Token grants'
      ],
      datasets: [{
        data: [tokenBalance, stakeBalance, withdrawalsTotal, grantBalance],
        backgroundColor: [
        '#505e5b',
        '#48dbb4',
        '#2f9278',
        '#FFCE56'
        ]
      }]
    };

    this.setState({
      yourAddress,
      tokenBalance,
      stakeBalance,
      grantBalance,
      grantStakeBalance,
      chartOptions,
      chartData,
      withdrawals,
      withdrawalsTotal
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
