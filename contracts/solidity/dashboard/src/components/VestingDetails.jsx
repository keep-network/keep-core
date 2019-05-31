import React, { Component } from 'react'
import { Table } from 'react-bootstrap'
import TableRow from './TableRow'
import moment from 'moment'

class VestingDetails extends Component {
  render() {
    if (this.props.details) {
      const { start, cliff, end, amount, released, grantedAmount, revocable, owner } = this.props.details
      const unreleased = grantedAmount ? grantedAmount - released : null
      return <div>
        <Table striped bordered condensed className="small table-sm">
          <tbody>
            <TableRow title="From">
              <a href="/">{owner}</a>
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
            
            <TableRow title="Granted Amount">
              { this.formatTokens(grantedAmount) }
            </TableRow>
            
            <TableRow title="Already released">
              { this.formatTokens(released) }
            </TableRow>
            
            <TableRow title="Unreleased">
              { this.formatTokens(unreleased) }
            </TableRow>
            <TableRow title="Revocable">
              { revocable }
            </TableRow>
          </tbody>
        </Table>
      </div>
    } else {
      return "Loading ..."
    }
  }

  formatDate(date) {
    if (! date) return
    const milliseconds = date * 1000
    return moment(milliseconds).format("MMMM Do YYYY, h:mm:ss a")
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
