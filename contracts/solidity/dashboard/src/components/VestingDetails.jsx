import React, { Component } from 'react'
import { Table } from 'react-bootstrap'
import TableRow from './TableRow'
import moment from 'moment'

class VestingDetails extends Component {
  render() {
    if (this.props.details) {
      const { start, cliff, end, amount, withdrawn, grantedAmount, revocable, grantManager, staked } = this.props.details
      const available = grantedAmount ? grantedAmount - withdrawn : null
      return <div>
        <Table striped bordered condensed className="small table-sm">
          <tbody>
            <TableRow title="Grant Manager">
              <a href="/">{grantManager}</a>
            </TableRow>

            <TableRow title="Start date">
              { this.formatDate(start) }
            </TableRow>

            <TableRow title="Cliff">
              { this.formatDate(cliff) }
            </TableRow>

            <TableRow title="End date">
              { this.formatDate(end) }
            </TableRow>

            <TableRow title="Total vesting">
              { this.formatTokens(amount) }
            </TableRow>

            <TableRow title="Vested">
              { this.formatTokens(grantedAmount) }
            </TableRow>

            <TableRow title="Withdrawn">
              { this.formatTokens(withdrawn) }
            </TableRow>

            <TableRow title="Available to withdraw">
              { this.formatTokens(available) }
            </TableRow>

            <TableRow title="Revocable">
              { revocable }
            </TableRow>

            <TableRow title="Staked">
              { this.formatTokens(staked) }
            </TableRow>
          </tbody>
        </Table>
      </div>
    } else {
      return 'Loading ...'
    }
  }

  formatDate(date) {
    if (! date) return
    const milliseconds = date * 1000
    return moment(milliseconds).format('MMMM Do YYYY, h:mm:ss a')
  }

  displayAmount(amount, decimals) {
    amount = amount / (10 ** decimals)
    return Math.round(amount * 10000) / 10000
  }

  formatTokens(amount) {
    if (amount == null) return
    const { decimals, symbol } = this.props.details
    const display = this.displayAmount(amount, decimals)
    return `${display} ${symbol}`
  }
}

export default VestingDetails
