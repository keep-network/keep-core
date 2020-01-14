import React, { Component } from 'react'
import { Table } from 'react-bootstrap'

class TokenGrantsTable extends Component {
  renderRow = (item, i) => (
    <tr key={'token-grant-'+i+'-from-'+item.grantManager}>
      <td>{item.formatted.amount}</td>
      <td><a href="/">{item.grantManager}</a></td>
    </tr>
  )

  render() {
    return (
      <Table className="small table-sm" condensed>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>Grant Manager</strong></th>
          </tr>
        </thead>
        <tbody>
          { this.pros.data && this.props.data.map(this.renderRow) }
        </tbody>
      </Table>
    )
  }
}

export default TokenGrantsTable
