import React, { Component } from 'react'
import { Table } from 'react-bootstrap'
import TokenGrantRevokeButton from './TokenGrantRevokeButton'

class TokenGrantManagerTable extends Component {
  renderRow = (item, i) => (
    <tr key={'token-grant-'+i+'-for-'+item.grantee}>
      <td>{item.formatted.amount}</td>
      <td><a href="/">{item.grantee}</a></td>
      <td><TokenGrantRevokeButton item={item} /></td>
    </tr>
  )

  render() {
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
          { this.props.data && this.props.data.map(this.renderRow) }
        </tbody>
      </Table>
    )
  }
}

export default TokenGrantManagerTable
