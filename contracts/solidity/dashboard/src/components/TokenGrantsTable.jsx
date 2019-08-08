import React, { Component } from 'react'
import { Table } from 'react-bootstrap'

class TokenGrantsTable extends Component {
  render() {
    if (this.props.data) {
      var rows = this.props.data.map(function(item, i){
        return (
          <tr key={"token-grant-"+i+"-from-"+item.grantManager}>
            <td>{item.formatted.amount}</td>
            <td><a href="/">{item.grantManager}</a></td>
          </tr>
        )
      })
    }
    return (
      <Table className="small table-sm" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>Grant Manager</strong></th>
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
