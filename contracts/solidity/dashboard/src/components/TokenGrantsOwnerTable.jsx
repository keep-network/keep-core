import React, { Component } from 'react'
import { Table } from 'react-bootstrap'

class TokenGrantsOwnerTable extends Component {
  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <tr key={"token-grant-"+i+"-for-"+item.beneficiary}>
            <td>{item.formatted.amount}</td>
            <td><a href="">{item.beneficiary}</a></td>
            <td><button>Revoke</button></td>
          </tr>
        )
      })
    }
    return (
      <Table className="small table-sm" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>To</strong></th>
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

export default TokenGrantsOwnerTable
