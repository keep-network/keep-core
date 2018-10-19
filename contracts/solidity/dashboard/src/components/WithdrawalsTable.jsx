import React, { Component } from 'react'
import { Table } from 'react-bootstrap'
import Withdrawal from './Withdrawal'

class WithdrawalsTable extends Component {

  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <Withdrawal key={i} withdrawal={item}/>
        )
      })
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
    )
  }
}

export default WithdrawalsTable
