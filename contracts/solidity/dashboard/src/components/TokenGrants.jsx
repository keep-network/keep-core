import React, { Component } from 'react'
import { Table } from 'react-bootstrap'

class TokenGrants extends Component {
  renderRow = (item, i) => (
    <tr key={'token-grant-'+i+'-from-'+item.grantManager} onClick={() => this.props.selectTokenGrant(i)}>
      <td>{item.id}</td>
      <td>{item.formatted.amount}</td>
      <td><a href="/">{item.grantManager}</a></td>
      <td>{item.formatted.start}</td>
      <td>{item.formatted.end}</td>
      <td>{item.formatted.cliff}</td>
      <td>{item.formatted.withdrawn}</td>
    </tr>
  )

  render() {
    return (
      <Table className="small table-sm" condensed hover>
        <thead>
          <tr>
            <th className="text-nowrap"><strong>Grant Id</strong></th>
            <th><strong>Amount</strong></th>
            <th><strong>Grant Manager</strong></th>
            <th><strong>Start</strong></th>
            <th><strong>End</strong></th>
            <th><strong>Cliff</strong></th>
            <th><strong>Released</strong></th>
          </tr>
        </thead>
        <tbody>
          { this.props.data && this.props.data.map(this.renderRow) }
        </tbody>
      </Table>
    )
  }
}

export default TokenGrants
