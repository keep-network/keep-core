import React from 'react'
import PendingUndelegationList from './PendingUndelegationList'

const PendingUndelegation = (props) => {
  return (
    <section id="pending-undelegation" className="tile">
      <h5>Pending Undelegation</h5>
      <div className="flex pending-undelegation-summary">
        <h2 className="balance flex flex-1">3,000 K</h2>
        <div className="flex flex-1 flex-column">
          <span className="text-label">UNDELEGATED ON</span>
          <span className="text-big">January 19, 2020</span>
        </div>
        <div className="flex flex-1 flex-column">
          <span className="text-label">UNDELEGATION PERIOD</span>
          <span className="text-big">1 week</span>
        </div>
      </div>
      <div>
        <PendingUndelegationList />
      </div>
    </section>
  )
}

export default PendingUndelegation
