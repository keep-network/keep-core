import React from 'react'
import { formatDate, displayAmount } from '../utils'

const PendingUndelegationList = ({ pendingUndelegations }) => {
  return (
    <div className="pending-undelegation-list">
      <div className="flex flex-1">
        <span className="text-label flex-3" >UNDELEGATION STARTED</span>
        <span className="text-label flex flex-1">AMOUNT (KEEP)</span>
      </div>
      <ul className="flex flex-column flex-1">
        {pendingUndelegations.map(renderUndelegationItem)}
      </ul>
    </div>
  )
}

const PendingUndelegationItem = ({ amount, createdAt }) => (
  <li className="pending-undelegation-item flex flex-row flex-1" >
    <span className="text-big flex-3">{formatDate(createdAt)}</span>
    <span className="text-big flex-1">{amount && `${displayAmount(amount)}`}</span>
  </li>
)

const renderUndelegationItem = (item) => <PendingUndelegationItem key={item.eventId} {...item}/>

export default PendingUndelegationList
