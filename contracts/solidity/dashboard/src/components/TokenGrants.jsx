import React, { Component } from 'react';
import { Table } from 'react-bootstrap';

class TokenGrants extends Component {

  render() {
    if (this.props.data) {
      var rows = this.props.data.map((item, i) => {
        return (
          <tr key={i} onClick={() => this.props.selectTokenGrant(i)}>
            <td>{item.formatted.amount}</td>
            <td><a href="">{item.owner}</a></td>
            <td>{item.formatted.start}</td>
            <td>{item.formatted.end}</td>
            <td>{item.formatted.cliff}</td>
            <td>{item.formatted.released}</td>
          </tr>
        );
      });
    }
    return (
      <Table className="small table-sm" condensed hover>
        <thead>
          <tr>
            <th><strong>Amount</strong></th>
            <th><strong>From</strong></th>
            <th><strong>Start</strong></th>
            <th><strong>End</strong></th>
            <th><strong>Cliff</strong></th>
            <th><strong>Released</strong></th>
          </tr>
        </thead>
        <tbody>
          { rows }
        </tbody>
      </Table>
    );
  }
}

export default TokenGrants;
