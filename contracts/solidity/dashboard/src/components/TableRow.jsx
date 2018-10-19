import React, { Component } from 'react'

class TableRow extends Component {
  render() {
    return (
      <tr>
        <th><strong>{ this.props.title }</strong></th>
        <td>
          { this.props.children }
        </td>
      </tr>
    )
  }
}

export default TableRow
