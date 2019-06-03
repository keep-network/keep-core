import React, { Component } from 'react'
import { Table } from 'react-bootstrap'

class TokenGrantsTable extends Component {
  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <tr key={"token-grant-"+i+"-from-"+item.owner}>
            <td>{item.formatted.amount}</td>
            <td><a href="/">{item.owner}</a></td>
          </tr>
        )
      })
    }
    return (
      <Table className="small table-sm" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>From</strong></th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </Table>
    )
  }
}

export default TokenGrantsTable
