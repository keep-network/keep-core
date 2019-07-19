import BigNumber from "bignumber.js"
import React, { Component } from 'react'
import { Pie } from 'react-chartjs-2'
import { Table, Col, Grid, Row, Tabs, Tab } from 'react-bootstrap'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom'
import moment from 'moment'
import { displayAmount, getWeb3 } from './utils'
import { Web3Context } from './components/WithWeb3Context'
import { getKeepToken, getTokenStaking, getTokenGrant } from './contracts'
import Header from './components/Header'
import Footer from './components/Footer'
import StakingForm from './components/StakingForm'
import StakingTable from './components/StakingTable'
import StakingDelegateForm from './components/StakingDelegateForm'
import StakingDelegateTokenGrantForm from './components/StakingDelegateTokenGrantForm'
import SigningForm from './components/SigningForm'
import WithdrawalsTable from './components/WithdrawalsTable'
import TokenGrantForm from './components/TokenGrantForm'
import TokenGrantsOwnerTable from './components/TokenGrantsOwnerTable'
import TokenGrants from './components/TokenGrants'
import VestingChart from './components/VestingChart'
import VestingDetails from './components/VestingDetails'
import TableRow from './components/TableRow'
import { colors } from './colors'

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
    this.state = {
      web3: {
        yourAddress: undefined,
        networkType: undefined,
        token: undefined,
        stakingContract: undefined,
        grantContract: undefined,
        utils: undefined,
        eth: undefined
      }
    }
    this.state.operatorChartData = {}
    this.state.stakeOwnerChartData = {}
  }

  async componentDidMount() {

    const web3 = await getWeb3()
    if (!web3) {
      this.setState({
        error: "No network detected. Do you have MetaMask installed?",
      })
      return
    }

    const contracts = await this.getContracts(web3)
    if (!contracts) {
      this.setState({
        error: "Failed to load contracts. Please check if Metamask is enabled and connected to the correct network.",
      })
      return
    }

    this.setState({
      web3: {
        yourAddress: (await web3.eth.getAccounts())[0],
        networkType: await web3.eth.net.getNetworkType(),
        token: contracts.token,
        stakingContract: contracts.stakingContract,
        grantContract: contracts.grantContract,
        utils: web3.utils,
        eth: web3.eth
      }
    })

    this.getData()
  }

  getContracts = async (web3) => {
    try {
      const token = await getKeepToken(web3, process.env.REACT_APP_TOKEN_ADDRESS)
      const stakingContract = await getTokenStaking(web3, process.env.REACT_APP_STAKING_ADDRESS)
      const grantContract = await getTokenGrant(web3, process.env.REACT_APP_TOKENGRANT_ADDRESS)
      return {
        token: token,
        stakingContract: stakingContract,
        grantContract: grantContract
      }
    } catch (e) {
      return null
    }
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
    const { web3, tokenBalance, operators, stakeBalance, grantBalance, grantStakeBalance,
      isTokenHolder, isOperator, isOperatorOfStakedTokenGrant, stakedGrant, stakeOwner, operatorChartData, stakeOwnerChartData, chartOptions, withdrawals, withdrawalsTotal, grantedToYou, grantedByYou,
      error } = this.state

    return (
      <Web3Context.Provider value={web3}>
        <div className="main">
          <Header networkType={web3.networkType}/>
          <Grid>
            <Row>
              <Col xs={12}>
                {error ?
                  <div className="alert alert-danger m-5" role="alert">{error}</div>:null
                }

                {isOperator && !isOperatorOfStakedTokenGrant ?
                  <div className="alert alert-info m-5" role="alert">You are registered as an operator for {stakeOwner}</div>:null
                }

                {isOperatorOfStakedTokenGrant ?
                  <div className="alert alert-info m-5" role="alert">
                    You are registered as a staked token grant operator for {stakedGrant.grantee} 
                  </div>:null
                }

                {!isTokenHolder && !isOperator ?
                  <div className="signing">
                    <div className="alert alert-info m-5" role="alert">Sorry, looks like you don't have any tokens to stake.</div>
                    <h3>Become an operator</h3>
                    <p>
                      To become an operator you must have a mutual agreement with the stake owner. This is achieved by creating
                      a signature of the stake owner address and sending it to the owner. Using the signature the owner can initiate
                      stake delegation and you will be able to participate in network operations on behalf of the stake owner.
                    </p>

                    <div className="signing-form well">
                      <SigningForm description="Sign stake owner address" defaultMessageToSign="0x0" />
                      <SigningForm
                        description="(Optional) Sign Token Grant contract address. This is required only for Token Grants stake operators"
                        defaultMessageToSign={ process.env.REACT_APP_TOKENGRANT_ADDRESS }/>
                    </div>

                  </div>:
                  <Tabs defaultActiveKey={1} id="dashboard-tabs">
                    <Tab eventKey={1} title="Overview">
                      {isOperator ?
                      <Row className="overview">
                        <Col xs={12} md={6}>
                          <Pie dataKey="name" data={ operatorChartData } options={ chartOptions } />
                        </Col>
                        <Col xs={12} md={6}>
                          <Table className="small table-sm" striped bordered condensed>
                            <tbody>
                              <TableRow title="Your wallet address">
                                { this.state.web3.yourAddress }
                              </TableRow>
                              <TableRow title="Tokens">
                                { displayAmount(tokenBalance, 18, 3) }
                              </TableRow>
                              <TableRow title="Staked">
                                { displayAmount(stakeBalance, 18, 3) }
                              </TableRow>
                              <TableRow title="Pending unstake">
                                { withdrawalsTotal }
                              </TableRow>
                            </tbody>
                          </Table>

                          {!isOperatorOfStakedTokenGrant ?
                            <StakingForm btnText="Unstake" action="unstake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>:
                            <div>
                              <StakingForm btnText="Unstake" action="unstake" stakingContractAddress={ process.env.REACT_APP_TOKENGRANT_ADDRESS }/>
                              <small>You can only unstake full amount. Partial unstake amounts are not yet supported.</small>
                            </div>
                          }

                        </Col>
                      </Row>:
                      <Row className="overview">
                        <Col xs={12} md={6}>
                          <Pie dataKey="name" data={ stakeOwnerChartData } options={ chartOptions } />
                        </Col>
                        <Col xs={12} md={6}>
                          <Table className="small table-sm" striped bordered condensed>
                            <tbody>
                              <TableRow title="Your wallet address">
                                { this.state.web3.yourAddress }
                              </TableRow>
                              <TableRow title="Tokens">
                                { displayAmount(tokenBalance, 18, 3) }
                              </TableRow>
                              <TableRow title="Staked">
                                { displayAmount(stakeBalance, 18, 3) }
                              </TableRow>
                              <TableRow title="Pending unstake">
                                { withdrawalsTotal }
                              </TableRow>
                              <TableRow title="Token Grants">
                                { displayAmount(grantBalance, 18, 3) }
                              </TableRow>
                              <TableRow title="Staked Token Grants">
                                { grantStakeBalance }
                              </TableRow>
                            </tbody>
                          </Table>
                        </Col>
                      </Row>
                      }
                      <Row>
                        {!isOperator ?
                        <Col sm={12}>
                          <h4>Delegated Stake</h4>
                          <StakingTable data={operators}/>
                        </Col>:null
                        }
                        <Col sm={12}>
                          <h4>Pending unstake</h4>
                          <WithdrawalsTable data={withdrawals}/>
                        </Col>
                      </Row>
                    </Tab>
                    {!isOperator ?
                      <Tab eventKey={2} title="Stake">
                        <Row>
                          <Col xs={12} md={12}>
                            <h3>Stake Delegation</h3>
                            <p>
                              Keep network does not require token owners to perform the day-to-day operations of staking 
                              with the private keys holding the tokens. This is achieved by stake delegation, where different
                              addresses hold different responsibilities and cold storage is supported to the highest extent practicable.
                            </p>
                            <StakingDelegateForm tokenBalance={tokenBalance} stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS } />
                            <hr></hr>
                            <h3>Stake Delegation (Simplified)</h3>
                            <p>
                              Simplified arrangement where you operate and receive rewards under one account.
                            </p>
                            <StakingForm btnText="Stake" action="stake" stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }/>
                          </Col>
                        </Row>
                      </Tab>:null
                    }
                    {!isOperator ?
                      <Tab eventKey={3} title="Token Grants">
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
                          <Col xs={12} md={12}>
                            <TokenGrants data={grantedToYou} selectTokenGrant={this.selectTokenGrant} />
                          </Col>
                        </Row>
                        <Row>
                          <Col xs={12} md={12}>
                            <h3>Stake Delegation of Token Grants</h3>
                            <p>
                              Keep network does not require token owners to perform the day-to-day operations of staking 
                              with the private keys holding the tokens. This is achieved by stake delegation, where different
                              addresses hold different responsibilities and cold storage is supported to the highest extent practicable.
                            </p>
                            
                            <StakingDelegateTokenGrantForm
                              tokenBalance={grantBalance}
                              tokenGrantContractAddress={ process.env.REACT_APP_TOKENGRANT_ADDRESS }
                              stakingContractAddress={ process.env.REACT_APP_STAKING_ADDRESS }
                            />
                          </Col>
                        </Row>
                      </Tab>:null
                    }
                    {!isOperator ?
                      <Tab eventKey={4} title="Create Token Grant">
                        <h3>Grant tokens</h3>
                        <p>You can grant tokens with a vesting schedule where balance released to the grantee
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
                      </Tab>:null
                    }
                  </Tabs>
                }
              </Col>
            </Row>
          </Grid>
          <Footer />
        </div>
      </Web3Context.Provider>
    )
  }

  async getData() {

    const token = this.state.web3.token
    const stakingContract = this.state.web3.stakingContract
    const grantContract = this.state.web3.grantContract
    const tokenBalance = await token.methods.balanceOf(this.state.web3.yourAddress).call()
    const stakeOwner = await stakingContract.methods.ownerOf(this.state.web3.yourAddress).call();
    const grantBalance = await grantContract.methods.balanceOf(this.state.web3.yourAddress).call()
    const grantStakeBalance = displayAmount(await grantContract.methods.stakeBalanceOf(this.state.web3.yourAddress).call(), 18, 3)

    let isTokenHolder = false;
    let isOperator = false;

    if (tokenBalance.gt(0)) {
      isTokenHolder = true;
    }

    if (stakeOwner !== "0x0000000000000000000000000000000000000000" && stakeOwner !== this.state.web3.yourAddress) {
      isOperator = true;
    }

    // Check if your account is an operator for a staked Token Grant.
    let stakedGrant
    let isOperatorOfStakedTokenGrant
    let stakedGrantByOperator = await grantContract.methods.grantStakes(this.state.web3.yourAddress).call()

    if (stakedGrantByOperator.stakingContract === stakingContract.address) {
      isOperatorOfStakedTokenGrant = true
      stakedGrant = await grantContract.methods.grants(stakedGrantByOperator.grantId.toString()).call()
    }

    // Calculate delegated stake balances
    let stakeBalance = await stakingContract.methods.balanceOf(this.state.web3.yourAddress).call()
    const operatorsAddresses = await stakingContract.methods.operatorsOf(this.state.web3.yourAddress).call()
    let operators = [];

    for(let i = 0; i < operatorsAddresses.length; i++) {
      let balance = await stakingContract.methods.balanceOf(operatorsAddresses[i]).call();
      if (!balance.isZero()) {
        let operator = {
          'address': operatorsAddresses[i],
          'amount': balance
        }
        operators.push(operator)
        stakeBalance = balance.add(stakeBalance)
      }
    }

    // Unstake withdrawals
    let withdrawalsByOperator = []

    if (isOperator) {
      withdrawalsByOperator.push(this.state.web3.yourAddress)
    } else {
      withdrawalsByOperator = operatorsAddresses;
    }

    const withdrawalDelay = (await stakingContract.methods.stakeWithdrawalDelay().call()).toNumber()
    let withdrawals = []
    let withdrawalsTotal = new BigNumber(0)

    for(let i=0; i < withdrawalsByOperator.length; i++) {
      const withdrawal = await stakingContract.methods.getWithdrawal(withdrawalsByOperator[i]).call()
      if (withdrawal[0] > 0) {
        const withdrawalAmount = displayAmount(withdrawal[0], 18, 3)
        withdrawalsTotal = withdrawalsTotal.plus(withdrawal[0])
        const availableAt = moment(withdrawal[1].toNumber()*1000).add(withdrawalDelay, 'seconds')
        let available = false
        const now = moment()
        if (availableAt.isSameOrBefore(now)) {
          available = true
        }

        withdrawals.push({
          'id': withdrawalsByOperator[i],
          'amount': withdrawalAmount,
          'available': available,
          'availableAt': availableAt.format("MMMM Do YYYY, h:mm:ss a")
          }
        )
      }
    }

    withdrawalsTotal = displayAmount(withdrawalsTotal, 18, 3)

    // Token Grants
    const grantIndexes = await grantContract.methods.getGrants(this.state.web3.yourAddress).call()
    let grantedToYou = []
    let grantedByYou = []

    for(let i=0; i < grantIndexes.length; i++) {
      const grant = await grantContract.methods.grants(grantIndexes[i].toNumber()).call()
      const grantedAmount = await grantContract.methods.grantedAmount(grantIndexes[i].toNumber()).call()
      const data = {
        'id': grantIndexes[i].toNumber(),
        'owner': this.state.web3.utils.toChecksumAddress(grant[0]),
        'grantee': this.state.web3.utils.toChecksumAddress(grant[1]),
        'revoked': grant[2],
        'revocable': grant[3],
        'amount': grant[4],
        'grantedAmount': grantedAmount,
        'end': grant[5].add(grant[6]),
        'start': grant[6],
        'cliff': grant[7],
        'withdrawn': grant[8],
        'staked': grant[9],
        'decimals': 18,
        'symbol': 'KEEP',
        'formatted': {
          'amount': displayAmount(grant[4], 18, 3),
          'end': moment((grant[5].add(grant[6])).mul(1000)).format("MMMM Do YYYY, h:mm:ss a"),
          'start': moment((grant[6].toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a"),
          'cliff': moment((grant[7].toNumber())* 1000).format("MMMM Do YYYY, h:mm:ss a"),
          'withdrawn': grant[8].toNumber()
        }
      }

      if (this.state.web3.yourAddress === data['owner']) {
        grantedByYou.push(data)
      } else if (this.state.web3.yourAddress === data['grantee']) {
        grantedToYou.push(data)
      }
    }

    let selectedGrantIndex = 0

    const chartOptions = {
      legend: {
        position: 'right'
      }
    }
    const operatorChartData = {
      labels: [
        'Delegated stake',
        'Pending unstake'
      ],
      datasets: [{
        data: [stakeBalance, withdrawalsTotal],
        backgroundColor: [
          colors.nandor,
          colors.turquoise
        ]
      }]
    }

    const stakeOwnerChartData = {
      labels: [
        'Tokens',
        'Delegated stake',
        'Pending unstake',
        'Token grants'
      ],
      datasets: [{
        data: [tokenBalance, stakeBalance, withdrawalsTotal, grantBalance],
        backgroundColor: [
          colors.nandor,
          colors.turquoise,
          colors.lochinvar,
          colors.goldenTainoi
        ]
      }]
    }

    this.setState({
      tokenBalance,
      operators,
      isTokenHolder,
      isOperator,
      isOperatorOfStakedTokenGrant,
      stakedGrant,
      stakeOwner,
      stakeBalance,
      grantBalance,
      grantStakeBalance,
      chartOptions,
      operatorChartData,
      stakeOwnerChartData,
      withdrawals,
      withdrawalsTotal,
      grantedToYou,
      grantedByYou,
      selectedGrantIndex
    })
  }
}

export default App
