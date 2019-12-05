import React from 'react'
import { Pie } from 'react-chartjs-2'
import { Table, Col, Row } from 'react-bootstrap'
import moment from 'moment'
import { displayAmount } from '../utils'
import StakingForm from './StakingForm'
import StakingTable from './StakingTable'
import WithdrawalsTable from './WithdrawalsTable'
import TableRow from './TableRow'
import { colors } from '../colors'
import WithWeb3Context from './WithWeb3Context'
import { withContractsDataContext } from './ContractsDataContextProvider'

class OverviewTab extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
          operators: [],
          stakeBalance: '',
          withdrawals: [],
          withdrawalsTotal: [],
          beneficiaryAddress: '',
          chartOptions: {
            legend: {
                position: 'right'
              }
          },
        }
    }

    componentDidMount() {
        this.getData();
    }

    componentDidUpdate(prevProps) {
        if(prevProps.web3.yourAddress !== this.props.web3.yourAddress)
          this.getData();
    }

    async getData() {
        const { web3: { token, stakingContract, grantContract, yourAddress, utils } } = this.props
        const { isOperator } = this.props
        if(!token.options.address || !stakingContract.options.address || !grantContract.options.address)
          return
    
        // Calculate delegated stake balances
        let stakeBalance = new utils.BN(await stakingContract.methods.balanceOf(yourAddress).call());
        const operatorsAddresses = await stakingContract.methods.operatorsOf(yourAddress).call()
        let operators = [];
    
        for(let i = 0; i < operatorsAddresses.length; i++) {
          let balance = new utils.BN(await stakingContract.methods.balanceOf(operatorsAddresses[i]).call())
          if (!balance.isZero()) {
            let operator = {
              'address': operatorsAddresses[i],
              'amount': balance.toString()
            }
            operators.push(operator)
            stakeBalance = balance.add(stakeBalance)
          }
        }
    
        // Unstake withdrawals
        let withdrawalsByOperator = isOperator ? [yourAddress] : operatorsAddresses
    
        const withdrawalDelay = await stakingContract.methods.stakeWithdrawalDelay().call();
        let withdrawals = []
        let withdrawalsTotal = new utils.BN(0);
    
        for(let i=0; i < withdrawalsByOperator.length; i++) {
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
              'availableAt': availableAt.format("MMMM Do YYYY, h:mm:ss a")
              }
            )
          }
        }
    
        this.setState({
          operators,
          stakeBalance,
          withdrawals,
          withdrawalsTotal: displayAmount(withdrawalsTotal, 18, 3),
          beneficiaryAddress: await this.getBeneficiaryAddress()
        })
    }

    getChartData = () => {
        const { stakeBalance, withdrawalsTotal } = this.state;
        const { isOperator, tokenBalance, grantBalance } = this.props;
        return isOperator ? 
        {
            labels: [
              'Delegated stake',
              'Pending unstake'
            ],
            datasets: [{
              data: [displayAmount(stakeBalance, 18, 3), withdrawalsTotal],
              backgroundColor: [
                colors.nandor,
                colors.turquoise
              ]
            }]
          } : 
          {
            labels: [
              'Tokens',
              'Delegated stake',
              'Pending unstake',
              'Token grants'
            ],
            datasets: [{
              data: [displayAmount(tokenBalance, 18, 3), displayAmount(stakeBalance, 18, 3), withdrawalsTotal, displayAmount(grantBalance, 18, 3)],
              backgroundColor: [
                colors.nandor,
                colors.turquoise,
                colors.lochinvar,
                colors.goldenTainoi
              ]
            }]
          }
    }

    getBeneficiaryAddress = async () => {
      const { web3: { utils, stakingContract, yourAddress }, isOperator } = this.props;
      const beneficiaryAddress = isOperator ? await stakingContract.methods.magpieOf(yourAddress).call() : ''
      return beneficiaryAddress && utils.toChecksumAddress(beneficiaryAddress)
    }

    render() {
        const { operators, chartOptions, withdrawals, withdrawalsTotal, stakeBalance, beneficiaryAddress } = this.state
        const { web3, isOperator, isOperatorOfStakedTokenGrant, tokenBalance, grantBalance, grantStakeBalance } = this.props;
        return (
            <>
              {isOperator ?
                <Row className="overview">
                  <Col xs={12} md={6}>
                    <Pie dataKey="name" data={this.getChartData()} options={ chartOptions } />
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
                            { withdrawalsTotal }
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
                <Row className="overview">
                  <Col xs={12} md={6}>
                    <Pie dataKey="name" data={this.getChartData()} options={ chartOptions } />
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
        );
    }
};

export default WithWeb3Context(withContractsDataContext(OverviewTab))