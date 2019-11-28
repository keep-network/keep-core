import React from 'react';
import { Pie } from 'react-chartjs-2'
import { Table, Col, Row } from 'react-bootstrap'
import moment from 'moment'
import { displayAmount } from '../utils'
import StakingForm from './StakingForm'
import StakingTable from './StakingTable'
import WithdrawalsTable from './WithdrawalsTable'
import TableRow from './TableRow'
import { colors } from '../colors'
import WithWeb3Context from './WithWeb3Context';

class OverviewTab extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            tokenBalance: '',
          operators: [],
          isTokenHolder: false,
          isOperator: false,
          isOperatorOfStakedTokenGrant: false,
          stakedGrant: '',
          stakeOwner: '',
          stakeBalance: '',
          grantBalance: '',
          grantStakeBalance: '',
          chartOptions: {
            legend: {
                position: 'right'
              }
          },
          operatorChartData: {},
          stakeOwnerChartData: {},
          withdrawals: [],
          withdrawalsTotal: [],
          grantedToYou: [],
          grantedByYou: [],
          selectedGrantIndex: 1
        }
    }

    async componentDidMount() {
        await this.getData();
    }

    componentDidUpdate(prevProps) {
        if(prevProps.web3.yourAddress !== this.props.web3.yourAddress)
          this.getData();
    }

    async getData() {
        const { web3: { token, stakingContract, grantContract, yourAddress, changeDefaultContract, utils } } = this.props;
        if(!token.options.address || !stakingContract.options.address || !grantContract.options.address)
          return
        // to display and to check is token holder  
        const tokenBalance = new utils.BN(await token.methods.balanceOf(yourAddress).call());
        //

        // check is operator
        const stakeOwner = await stakingContract.methods.ownerOf(yourAddress).call();
        //

        // To display only
        const grantBalance = await grantContract.methods.balanceOf(yourAddress).call()
        const grantStakeBalance = displayAmount(await grantContract.methods.stakeBalanceOf(yourAddress).call(), 18, 3)
        //

        let isTokenHolder = tokenBalance.gt(new utils.BN(0));
        let isOperator = stakeOwner !== "0x0000000000000000000000000000000000000000" && utils.toChecksumAddress(yourAddress) !== utils.toChecksumAddress(stakeOwner)
    
        // Check if your account is an operator for a staked Token Grant.
        let stakedGrant
        let isOperatorOfStakedTokenGrant
        let stakedGrantByOperator = await grantContract.methods.grantStakes(yourAddress).call()
    
        if (stakedGrantByOperator.stakingContract === stakingContract.options.address) {
          isOperatorOfStakedTokenGrant = true
          stakedGrant = await grantContract.methods.grants(stakedGrantByOperator.grantId.toString()).call()
          changeDefaultContract(grantContract);
        }
    
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
        let withdrawalsByOperator = []
    
        if (isOperator) {
          withdrawalsByOperator.push(yourAddress)
        } else {
          withdrawalsByOperator = operatorsAddresses;
        }
    
        const withdrawalDelay = await stakingContract.methods.stakeWithdrawalDelay().call();
        let withdrawals = []
        let withdrawalsTotal = new utils.BN(0);
    
        for(let i=0; i < withdrawalsByOperator.length; i++) {
            console.log('withdrawal by operator', withdrawalsByOperator[i])
          const withdrawal = await stakingContract.methods.getWithdrawal(withdrawalsByOperator[i]).call()
          console.log('with', withdrawal)
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
    
        withdrawalsTotal = displayAmount(withdrawalsTotal, 18, 3);
    
        this.setState({
          tokenBalance,
          operators,
          isTokenHolder,
          isOperator,
          isOperatorOfStakedTokenGrant,
          stakeBalance,
          grantBalance,
          grantStakeBalance,
          withdrawals,
          withdrawalsTotal,
        })
    }

    getChartData = () => {
        const { isOperator, tokenBalance, stakeBalance, withdrawalsTotal, grantBalance } = this.state;
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

    render() {
        const { tokenBalance, operators, stakeBalance, grantBalance, grantStakeBalance,isOperator, isOperatorOfStakedTokenGrant, chartOptions, withdrawals, withdrawalsTotal } = this.state
      
        const { web3 } = this.props;
        const { error } = web3;
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
                        <StakingForm btnText="Unstake" action="unstake" />:
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
            </>
        );
    }
};

export default WithWeb3Context(OverviewTab)