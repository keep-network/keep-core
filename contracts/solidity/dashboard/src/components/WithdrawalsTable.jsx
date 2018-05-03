import React, { Component } from 'react';
import { Table, Button } from 'react-bootstrap';
import { getTokenStaking } from '../contracts';
import Network from '../network';

class Withdrawal extends Component {

  async finishUnstake(withdrawalId) {
    const accounts = await Network.getAccounts();
    const stakingContract = await getTokenStaking(process.env.REACT_APP_STAKING_ADDRESS);
    stakingContract.finishUnstake(withdrawalId, {from: accounts[0], gas: 110000});
  }

  render() {

    let action = 'N/A';
    if (this.props.withdrawal.available) {
      action = <Button bsSize="small" bsStyle="primary" onClick={()=>this.finishUnstake(this.props.withdrawal.id)}>Finish Unstake</Button>;
    }

    return (
      <tr>
        <td>{this.props.withdrawal.amount}</td>
        <td className="text-mute">{this.props.withdrawal.availableAt}</td>
        <td>{action}</td>
      </tr>
    );
  }
}

class WithdrawalsTable extends Component {

  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <Withdrawal key={i} withdrawal={item}/>
        );
      });
    }

    return (
      <Table className="small table-sm" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>Available At</strong></th>
            <th><strong>Action</strong></th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </Table>
    );
  }
}

export default WithdrawalsTable;
