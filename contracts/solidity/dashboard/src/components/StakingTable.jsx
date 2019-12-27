import React from 'react'
import { Table } from 'react-bootstrap'
import { displayAmount } from '../utils'
import withWeb3Context from './WithWeb3Context'
import UndelegateStakeButton from './UndelegateStakeButton'

const StakingTable = (props) => (
  <Table className="small table-sm" condensed>
    <thead>
      <tr>
        <th><strong>Amount</strong></th>
        <th><strong>Operator</strong></th>
        <th><strong>Action</strong></th>
      </tr>
    </thead>
    <tbody>
      { props.data && props.data.map(renderRow) }
    </tbody>
  </Table>
)

const renderRow = (item, i) => (
  <tr key={`stake-delegate-${i}-to-${item.address}`}>
    <td>{displayAmount(item.amount, 18, 3)}</td>
    <td><a href="/">{item.address}</a></td>
    <td><UndelegateStakeButton amount={item.amount} operator={item.address}/></td>
  </tr>
)

export default withWeb3Context(StakingTable)
