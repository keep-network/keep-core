import React, { Component } from 'react';
import { Table } from 'react-bootstrap';

class WithdrawalsTable extends Component {
  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <tr key={i}>
            <td>{item.amount}</td>
            <td>{item.availableAt}</td>
          </tr>
        );
      });
    }
    return (
      <Table className="small" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>Available</strong></th>
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
