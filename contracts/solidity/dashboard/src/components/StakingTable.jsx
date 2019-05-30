import React from 'react'
import { Table } from 'react-bootstrap'
import { displayAmount } from '../utils'
import WithWeb3Context from './WithWeb3Context'
import UndelegateStakeButton from './UndelegateStakeButton'

const StakingTable = (props) => {
  if (props.data) {
    var rows = props.data.map(function(item, i){
      return (
        <tr key={`stake-delegate-${i}-to-${item.address}`}>
          <td>{displayAmount(item.amount,18, 3)}</td>
          <td><a href="/">{item.address}</a></td>
          <td><UndelegateStakeButton amount={item.amount} operator={item.address}/></td>
        </tr>
      )
    })
  }
  return (
    <Table className="small table-sm" condensed>
      <thead>
        <tr>
          <th><strong>Amount</strong></th>
          <th><strong>Operator</strong></th>
          <th><strong>Action</strong></th>
        </tr>
      </thead>
      <tbody>
        { rows }
      </tbody>
    </Table>
  )
}

export default WithWeb3Context(StakingTable)
