import React from 'react'
import { Pie } from 'react-chartjs-2'
import { Table, Col, Row } from 'react-bootstrap'
import moment from 'moment'
import { displayAmount, formatAmount } from '../utils'
import StakingForm from './StakingForm'
import StakingTable from './StakingTable'
import WithdrawalsTable from './WithdrawalsTable'
import TableRow from './TableRow'
import { colors } from '../colors'
import withWeb3Context from './WithWeb3Context'
import { withContractsDataContext } from './ContractsDataContextProvider'

class OverviewTab extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      operators: [],
      stakeBalance: 0,
      withdrawals: [],
      withdrawalsTotal: 0,
      beneficiaryAddress: '',
      chartOptions: {
        legend: {
          position: 'right',
        },
      },
      shouldSubscribeToEvent: true,
    }
  }

  componentDidMount() {
    this.getData()
    this.subscribeToEvent()
  }

  componentDidUpdate(prevProps, prevState) {
    if (prevProps.web3.yourAddress !== this.props.web3.yourAddress) {
      this.getData()
    }
    if (!prevState.shouldSubscribeToEvent && this.state.shouldSubscribeToEvent) {
      this.subscribeToEvent()
    }
  }

    subscribeToEvent = () => {
      const { web3EventProvider: { eventStakingContract } } = this.props
      this.setState({ shouldSubscribeToEvent: false })
      eventStakingContract.once('InitiatedUnstake', this.subscribeEvent)
    }

    subscribeEvent = async (error, event) => {
      const { returnValues: { value, operator, createdAt } } = event
      const { web3: { utils, stakingContract } } = this.props

      const withdrawalDelay = await stakingContract.methods.stakeWithdrawalDelay().call()
      const availableAt = moment(createdAt * 1000).add(withdrawalDelay, 'seconds')
      const withdrawal = {
        'id': operator,
        'amount': displayAmount(value, 18, 3),
        'available': availableAt.isSameOrBefore(moment()),
        'availableAt': availableAt.format('MMMM Do YYYY, h:mm:ss a'),
      }
      const withdrawals = [...this.state.withdrawals, withdrawal]
      const withdrawalsTotal = new utils.BN(this.state.withdrawalsTotal).add(utils.toBN(value))
      const stakeBalance = this.state.stakeBalance.sub(utils.toBN(value))
      const operators = this.state.operators.filter(({ address }) => address !== operator)

      this.setState({
        stakeBalance,
        operators,
        withdrawals,
        withdrawalsTotal,
        shouldSubscribeToEvent: true,
      })
    }

    async getData() {
      const { web3: { token, stakingContract, grantContract, yourAddress, utils } } = this.props
      const { isOperator } = this.props
      if (!token.options.address || !stakingContract.options.address || !grantContract.options.address) {
        return
      }

      // Calculate delegated stake balances
      let stakeBalance = new utils.BN(isOperator ? await stakingContract.methods.balanceOf(yourAddress).call(): 0)
      const operatorsAddresses = await stakingContract.methods.operatorsOf(yourAddress).call()
      const operators = []

      for (let i = 0; i < operatorsAddresses.length; i++) {
        const balance = new utils.BN(await stakingContract.methods.balanceOf(operatorsAddresses[i]).call())
        if (!balance.isZero()) {
          const operator = {
            'address': operatorsAddresses[i],
            'amount': balance.toString(),
          }
          operators.push(operator)
          stakeBalance = balance.add(stakeBalance)
        }
      }

      // Unstake withdrawals
      const withdrawalsByOperator = isOperator ? [yourAddress] : operatorsAddresses

      const withdrawalDelay = await stakingContract.methods.stakeWithdrawalDelay().call()
      const withdrawals = []
      let withdrawalsTotal = new utils.BN(0)

      for (let i=0; i < withdrawalsByOperator.length; i++) {
        const withdrawal = await stakingContract.methods.getWithdrawal(withdrawalsByOperator[i]).call()
        if (withdrawal.amount > 0) {
          withdrawalsTotal = withdrawalsTotal.add(new utils.BN(withdrawal.amount))

          const availableAt = moment(withdrawal.createdAt * 1000).add(withdrawalDelay, 'seconds')
          let available = false
          const now = moment()
          if (availableAt.isSameOrBefore(now)) {
            available = true
          }

          withdrawals.push({
            'id': withdrawalsByOperator[i],
            'amount': displayAmount(withdrawal.amount, 18, 3),
            'available': available,
            'availableAt': availableAt.format('MMMM Do YYYY, h:mm:ss a'),
          },
          )
        }
      }

      this.setState({
        operators,
        stakeBalance,
        withdrawals,
        withdrawalsTotal,
        beneficiaryAddress: await this.getBeneficiaryAddress(),
      })
    }

    getChartData = () => {
      const { stakeBalance, withdrawalsTotal } = this.state
      const { isOperator, tokenBalance, grantBalance } = this.props
      return isOperator ?
        {
          labels: [
            'Delegated stake',
            'Pending unstake',
          ],
          datasets: [{
            data: [displayAmount(stakeBalance, 18, 3), displayAmount(withdrawalsTotal, 18, 3)],
            backgroundColor: [
              colors.nandor,
              colors.turquoise,
            ],
          }],
        } :
        {
          labels: [
            'Tokens',
            'Delegated stake',
            'Pending unstake',
            'Token grants',
          ],
          datasets: [{
            data: [displayAmount(tokenBalance, 18, 3), displayAmount(stakeBalance, 18, 3), displayAmount(withdrawalsTotal, 18, 3), displayAmount(grantBalance, 18, 3)],
            backgroundColor: [
              colors.nandor,
              colors.turquoise,
              colors.lochinvar,
              colors.goldenTainoi,
            ],
          }],
        }
    }

    renderChart = () => {
      const { web3: { utils } } = this.props
      const chartData = this.getChartData()
      const shouldRenderChart = chartData.datasets[0].data.some((value) => !utils.toBN(formatAmount(value || 0, 18)).isZero())

      return ( shouldRenderChart ?
        <Pie dataKey="name" data={chartData} options={this.state.chartOptions} /> :
        <div className="alert alert-info m-5" role="alert">It looks like You don&apos;t have any tokens or delegated stake.</div>
      )
    }

    getBeneficiaryAddress = async () => {
      const { web3: { utils, stakingContract, yourAddress }, isOperator } = this.props
      const beneficiaryAddress = isOperator ? await stakingContract.methods.magpieOf(yourAddress).call() : ''
      return beneficiaryAddress && utils.toChecksumAddress(beneficiaryAddress)
    }

    render() {
      const { operators, withdrawals, withdrawalsTotal, stakeBalance, beneficiaryAddress } = this.state
      const { web3, isOperator, isOperatorOfStakedTokenGrant, tokenBalance, grantBalance, grantStakeBalance } = this.props
      return (
        <>
              {isOperator ?
                <Row>
                  <Col xs={12} md={6}>
                    {this.renderChart()}
                  </Col>
                  <Col xs={12} md={6}>
                    <Table className="small table-sm" striped bordered condensed>
                      <tbody>
                        <TableRow title="Your wallet address">
                          { web3.yourAddress }
                        </TableRow>
                        <TableRow title="Beneficiary address">
                          { beneficiaryAddress }
                        </TableRow>
                        <TableRow title="Tokens">
                          { displayAmount(tokenBalance, 18, 3) }
                        </TableRow>
                        <TableRow title="Staked">
                          { displayAmount(stakeBalance, 18, 3) }
                        </TableRow>
                        <TableRow title="Pending unstake">
                          { displayAmount(withdrawalsTotal, 18, 3) }
                        </TableRow>
                      </tbody>
                    </Table>
                    {!isOperatorOfStakedTokenGrant
                      ? <StakingForm btnText="Unstake" action="unstake" />
                      :
                      <div>
                        <StakingForm btnText="Unstake" action="unstake" />
                        <small>You can only unstake full amount. Partial unstake amounts are not yet supported.</small>
                      </div>
                    }
                  </Col>
                </Row>:
                <Row>
                  <Col xs={12} md={6}>
                    {this.renderChart()}
                  </Col>
                  <Col xs={12} md={6}>
                    <Table className="small table-sm" striped bordered condensed>
                      <tbody>
                        <TableRow title="Your wallet address">
                          { web3.yourAddress }
                        </TableRow>
                        <TableRow title="Tokens">
                          { displayAmount(tokenBalance, 18, 3) }
                        </TableRow>
                        <TableRow title="Staked">
                          { displayAmount(stakeBalance, 18, 3) }
                        </TableRow>
                        <TableRow title="Pending unstake">
                          { displayAmount(withdrawalsTotal, 18, 3) }
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
                {!isOperator &&
                  <Col sm={12}>
                    <h4>Delegated Stake</h4>
                    <StakingTable data={operators}/>
                  </Col>
                }
                <Col sm={12}>
                  <h4>Pending unstake</h4>
                  <WithdrawalsTable data={withdrawals}/>
                </Col>
              </Row>
        </>
      )
    }
};

export default withWeb3Context(withContractsDataContext(OverviewTab))
