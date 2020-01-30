import React, { useEffect, useContext, useRef } from 'react'
import PendingUndelegationList from './PendingUndelegationList'
import { useFetchData } from '../../hooks/useFetchData'
import { operatorService } from './service'
import { formatDate, displayAmount } from '../../utils'
import { LoadingOverlay } from '../Loadable'
import { Web3Context } from '../WithWeb3Context'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = (props) => {
  const event = useRef(null)
  const state = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: { pendingUnstakeBalance, undelegatedOn, stakeWithdrawalDelay, pendinUndelegations } } = state
  const { stakingContract, yourAddress } = useContext(Web3Context)

  useEffect(() => {
    event.current = stakingContract.events.InitiatedUnstake({ filter: { operator: yourAddress } }, (error, event) => {
      console.log('subscribed to event', error, event )
    })

    return () => {
      console.log('unmount')
      event.current.unsubscribe((error, suscces) => console.log('unsub', error, suscces))
    }
  }, [])

  return (
    <section id="pending-undelegation" className="tile">
      <LoadingOverlay isFetching={isFetching}>
        <h5>Pending Undelegation</h5>
        <div className="flex pending-undelegation-summary">
          <h2 className="balance flex flex-1">{pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)} K`}</h2>
          <div className="flex flex-1 flex-column">
            <span className="text-label">UNDELEGATED ON</span>
            <span className="text-big">{formatDate(undelegatedOn)}</span>
          </div>
          <div className="flex flex-1 flex-column">
            <span className="text-label">UNDELEGATION PERIOD</span>
            <span className="text-big">{stakeWithdrawalDelay}</span>
          </div>
        </div>
        <div>
          <PendingUndelegationList pendingUndelegations={pendinUndelegations} />
        </div>
      </LoadingOverlay>
    </section>
  )
}

export default PendingUndelegation
