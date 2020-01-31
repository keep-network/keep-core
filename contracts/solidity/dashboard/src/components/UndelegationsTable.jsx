import React, { Component } from 'react'
import { Table } from 'react-bootstrap'
import Undelegation from './Undelegation'

class UndelegationsTable extends Component {
  renderRow = (item, index) => {
    return <Undelegation key={index} undelegation={item}/>
  }

  render() {
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
          { this.props.data && this.props.data.map(this.renderRow) }
        </tbody>
      </Table>
    )
  }
}

export default UndelegationsTable
