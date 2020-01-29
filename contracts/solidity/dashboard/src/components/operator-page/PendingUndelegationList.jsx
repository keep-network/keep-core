import React from 'react'

const PendingUndelegationList = (props) => {
  return (
    <>
        <div className="flex flex-1">
          <span className="text-label flex-2" >UNDELEGATION PROCESS INITIATED</span>
          <span className="text-label flex flex-1">AMOUNT</span>
        </div>
        <ul className="pending-undelegation-list flex flex-column flex-1">
          <li className="pending-undelegation-item flex flex-row flex-1" >
            <span className="text-big flex-2">January 12, 2020</span>
            <span className="text-big flex-1">1,000 K</span>
          </li>
          <li className="pending-undelegation-item flex flex-row flex-1" >
            <span className="text-big flex-2">January 10, 2020</span>
            <span className="text-big flex-1">1,000 K</span>
          </li>
          <li className="pending-undelegation-item flex flex-row flex-1" >
            <span className="text-big flex-2">January 9, 2020</span>
            <span className="text-big flex-1">1,000 K</span>
          </li>
        </ul>
    </>
  )
}

export default PendingUndelegationList
