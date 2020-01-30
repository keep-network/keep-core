import React from 'react'
import { formatDate, displayAmount } from '../utils'

const PendingUndelegationList = ({ pendingUndelegations }) => {
  return (
    <>
        <div className="flex flex-1">
          <span className="text-label flex-2" >UNDELEGATION PROCESS INITIATED</span>
          <span className="text-label flex flex-1">AMOUNT</span>
        </div>
        <ul className="pending-undelegation-list flex flex-column flex-1">
          {pendingUndelegations.map(renderUndelegationItem)}
        </ul>
    </>
  )
}

const PendingUndelegationItem = ({ amount, createdAt }) => (
  <li className="pending-undelegation-item flex flex-row flex-1" >
    <span className="text-big flex-2">{formatDate(createdAt)}</span>
    <span className="text-big flex-1">{amount && `${displayAmount(amount)} K`}</span>
  </li>
)

const renderUndelegationItem = (item) => <PendingUndelegationItem key={item.eventId} {...item}/>

export default PendingUndelegationList
