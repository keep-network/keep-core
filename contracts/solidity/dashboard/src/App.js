import React, { Component } from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import moment from 'moment';
import { displayAmount } from './utils';
import Network from './network';
import { getKeepToken, getTokenStaking, getTokenGrant } from './contracts';
import Header from './components/Header';
import Footer from './components/Footer';
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
  }

  componentDidMount() {
    this.getData();
  }

  render() {
    const { yourAddress, tokenBalance, stakeBalance, withdrawals, withdrawalsTotal,
      totalAvailableToStake, totalAvailableToUnstake } = this.state;

    return (
      <div className="main">
      <Header />
      <Grid>
        <Row>
          <Col xs={12}>
            <Tabs defaultActiveKey={1} id="dashboard-tabs">
              <Tab eventKey={1} title="Overview">
                <Row className="overview">
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
                        <TableRow title="Pending unstake">
                          { withdrawalsTotal }
                        </TableRow>
                      </tbody>
                    </Table>
                    <h6>You can stake up to { totalAvailableToStake } KEEP</h6>
                    <StakingForm btnText="Stake" action="stake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>
                    <h6>You can unstake up to { totalAvailableToUnstake } KEEP</h6>
                    <StakingForm btnText="Unstake" action="unstake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>
                  </Col>
                </Row>
                <Row>
                    <Col xs={12} md={6}>
                      <h4>Pending unstake</h4>
                      <WithdrawalsTable data={withdrawals}/>
                    </Col>
                  </Row>
              </Tab>
            </Tabs>
          </Col>
        </Row>
      </Grid>
      <Footer />
      </div>
    )
  }

  async getData() {

    // Your address
    const accounts = await Network.getAccounts();
    const yourAddress = accounts[0];

    // Contracts
    const token = await getKeepToken(process.env.REACT_APP_TOKEN_ADDRESS);
    const stakingContract = await getTokenStaking(process.env.REACT_APP_STAKING_ADDRESS);

    // Balances
    const tokenBalance = displayAmount(await token.balanceOf(yourAddress), 18, 3);
    const stakeBalance = displayAmount(await stakingContract.stakeBalanceOf(yourAddress), 18, 3);
    const totalAvailableToStake = parseInt(tokenBalance)+parseInt(grantBalance);
    const totalAvailableToUnstake = parseInt(stakeBalance)+parseInt(grantStakeBalance);

    // Unstake withdrawals
    const withdrawalIndexes = await stakingContract.getWithdrawals(yourAddress);
    const withdrawalDelay = (await stakingContract.withdrawalDelay()).toNumber()
    let withdrawals = [];
    let withdrawalsTotal = 0;

    for(let i=0; i < withdrawalIndexes.length; i++) {
      const withdrawalId = withdrawalIndexes[i].toNumber();
      const withdrawal = await stakingContract.getWithdrawal(withdrawalId);
      const withdrawalAmount = displayAmount(withdrawal[1], 18, 3);
      withdrawalsTotal += withdrawal[1].toNumber();
      const availableAt = moment(withdrawal[2].toNumber()*1000).add(withdrawalDelay, 'seconds');
      let available = false;
      let now = moment();
      if (now > availableAt) {
        available = true;
      }

      withdrawals.push({
        'id': withdrawalId,
        'amount': withdrawalAmount,
        'available': available,
        'availableAt': availableAt.format("MMMM Do YYYY, h:mm:ss a")
        }
      );
    }

    withdrawalsTotal = displayAmount(withdrawalsTotal, 18, 3);

    this.setState({
      yourAddress,
      tokenBalance,
      stakeBalance,
      withdrawals,
      withdrawalsTotal,
      totalAvailableToStake,
      totalAvailableToUnstake
    })
  }
}
export default App;
