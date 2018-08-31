import Web3 from 'web3'
import BigNumber from "bignumber.js"
import React, { Component } from 'react'
import { Pie } from 'react-chartjs-2'
import { Table, Col, Grid, Row, Tabs, Tab } from 'react-bootstrap'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import moment from 'moment'
import { displayAmount } from './utils'
import Network from './network'
import { getKeepToken, getTokenStaking, getTokenGrant } from './contracts'
import Header from './components/Header'
import Footer from './components/Footer'
import StakingForm from './components/StakingForm'
import WithdrawalsTable from './components/WithdrawalsTable'
import TokenGrantsTable from './components/TokenGrantsTable'
import TokenGrantForm from './components/TokenGrantForm'
import TokenGrantsOwnerTable from './components/TokenGrantsOwnerTable'
import TokenGrants from './components/TokenGrants'
import VestingChart from './components/VestingChart'
import VestingDetails from './components/VestingDetails'

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
    this.state.chartData = {}
  }

  componentDidMount() {

    Network.getNetworkType().then(result => {
      this.setState({
        networkType: result,
      })
    })
  
    this.getData().catch(error => {
      this.setState({
        error: error,
      })
    })
  }

  selectTokenGrant = (i) => {
    this.setState(
      { selectedGrantIndex: i }
    )
  }

  selectedGrant = () => {
    if (this.state.grantedToYou) {
      return this.state.grantedToYou[this.state.selectedGrantIndex]
    } else {
      return {}
    }
  }

  render() {
    const { yourAddress, tokenBalance, stakeBalance, grantBalance, grantStakeBalance, 
      chartOptions, chartData, withdrawals, withdrawalsTotal, grantedToYou, grantedByYou,
      totalAvailableToStake, totalAvailableToUnstake, error } = this.state

    return (
      <div className="main">
        <Header networkType={this.state.networkType}/>
        <Grid>
          <Row>
            <Col xs={12}>
              {this.state.error ?
                <div className="alert alert-danger m-5" role="alert">{error}</div>:null
              }
              <Tabs defaultActiveKey={1} id="dashboard-tabs">
                <Tab eventKey={1} title="Overview">
                  <Row className="overview">
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
                          <TableRow title="Pending unstake">
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
                    <Col xs={12} md={6}>
                      <h4>Tokens granted to you</h4>
                      <TokenGrantsTable data={ grantedToYou }/>
                    </Col>
                  </Row>
                </Tab>
                <Tab eventKey={2} title="Token Grants">
                  <h3>Tokens granted to you</h3>
                  <Row>
                    <Col xs={12} md={6}>
                      <VestingDetails
                        details={this.selectedGrant()}
                      />
                    </Col>
                    <Col xs={12} md={6}>
                      <VestingChart details={this.selectedGrant()}/>
                    </Col>
                  </Row>
                  <Row>
                    <TokenGrants data={grantedToYou} selectTokenGrant={this.selectTokenGrant} />
                  </Row>
                </Tab>
                <Tab eventKey={3} title="Create Token Grant">
                  <h3>Grant tokens</h3>
                  <p>You can grant tokens with a vesting schedule where balance released to the beneficiary 
                    gradually in a linear fashion until start + duration. By then all of the balance will have vested.
                    You must approve the amount you want to grant by calling approve() method of the token contract first
                  </p>
                  <Row>
                    <Col xs={12} md={8}>
                      <TokenGrantForm tokenGrantContractAddress={ process.env.REACT_APP_TOKENGRANT_ADDRESS }/>
                    </Col>
                  </Row>
                  <Row>
                    <h3>Granted by you</h3>
                    <Col xs={12}>
                      <TokenGrantsOwnerTable data={ grantedByYou }/>
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
    const accounts = await Network.getAccounts()
    const yourAddress = accounts[0]

    // Contracts
    let token, stakingContract, grantContract
    try {
      token = await getKeepToken(process.env.REACT_APP_TOKEN_ADDRESS)
      stakingContract = await getTokenStaking(process.env.REACT_APP_STAKING_ADDRESS)
      grantContract = await getTokenGrant(process.env.REACT_APP_TOKENGRANT_ADDRESS)
    } catch (e) {
      throw "Failed to load contracts. Please check if Metamask is enabled and connected to the correct network."
    }

    // Balances
    const tokenBalance = displayAmount(await token.balanceOf(yourAddress), 18, 3)
    const stakeBalance = displayAmount(await stakingContract.stakeBalanceOf(yourAddress), 18, 3)
    const grantBalance = displayAmount(await grantContract.balanceOf(yourAddress), 18, 3)
    const grantStakeBalance = displayAmount(await grantContract.stakeBalanceOf(yourAddress), 18, 3)
    const totalAvailableToStake = parseInt(tokenBalance)+parseInt(grantBalance)
    const totalAvailableToUnstake = parseInt(stakeBalance)+parseInt(grantStakeBalance)

    // Unstake withdrawals
    const withdrawalIndexes = await stakingContract.getWithdrawals(yourAddress)
    const withdrawalDelay = (await stakingContract.withdrawalDelay()).toNumber()
    let withdrawals = []
    let withdrawalsTotal = new BigNumber(0)

    for(let i=0; i < withdrawalIndexes.length; i++) {
      const withdrawalId = withdrawalIndexes[i].toNumber()
      const withdrawal = await stakingContract.getWithdrawal(withdrawalId)
      const withdrawalAmount = displayAmount(withdrawal[1], 18, 3)
      withdrawalsTotal = withdrawalsTotal.plus(withdrawal[1])
      const availableAt = moment(withdrawal[2].toNumber()*1000).add(withdrawalDelay, 'seconds')
      let available = false
      const now = moment()
      if (availableAt.isSameOrBefore(now)) {
        available = true
      }

      withdrawals.push({
        'id': withdrawalId,
        'amount': withdrawalAmount,
        'available': available,
        'availableAt': availableAt.format("MMMM Do YYYY, h:mm:ss a")
        }
      )
    }

    withdrawalsTotal = displayAmount(withdrawalsTotal, 18, 3)

    // Token Grants
    const grantIndexes = await grantContract.getGrants(yourAddress)
    let grantedToYou = []
    let grantedByYou = []

    for(let i=0; i < grantIndexes.length; i++) {
      const grant = await grantContract.grants(grantIndexes[i].toNumber())
      const grantedAmount = await grantContract.grantedAmount(grantIndexes[i].toNumber())
      const data = {
        'owner': Web3.utils.toChecksumAddress(grant[0]),
        'beneficiary': Web3.utils.toChecksumAddress(grant[1]),
        'locked': grant[2],
        'revoked': grant[3],
        'revocable': grant[4],
        'amount': grant[5],
        'grantedAmount': grantedAmount,
        'end': grant[6].plus(grant[7]),
        'start': grant[7],
        'cliff': grant[8],
        'released': grant[9],
        'decimals': 18,
        'symbol': 'KEEP',
        'formatted': {
          'amount': displayAmount(grant[5].toNumber(), 18, 3),
          'end': moment(((grant[6].plus(grant[7])).toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a"),
          'start': moment((grant[7].toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a"),
          'cliff': moment((grant[8].toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a"),
          'released': grant[9].toNumber()
        }
      }

      if (yourAddress === data['owner']) {
        grantedByYou.push(data)
      } else if (yourAddress === data['beneficiary']) {
        grantedToYou.push(data)
      }
    }

    let selectedGrantIndex = 0

    const chartOptions = {
      legend: {
        position: 'right'
      }
    }
    const chartData = {
      labels: [
        'Tokens',
        'Staked',
        'Pending unstake',
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
    }

    this.setState({
      yourAddress,
      tokenBalance,
      stakeBalance,
      grantBalance,
      grantStakeBalance,
      chartOptions,
      chartData,
      withdrawals,
      withdrawalsTotal,
      grantedToYou,
      grantedByYou,
      selectedGrantIndex,
      totalAvailableToStake,
      totalAvailableToUnstake
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

export default App
