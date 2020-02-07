import React, { useEffect } from 'react'
import PendingUndelegationList from './PendingUndelegationList'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import web3Utils from 'web3-utils'
import { displayAmount } from '../utils'
import { LoadingOverlay } from './Loadable'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
    pendinUndelegations,
  } } = state

  useEffect(() => {
    if (latestUnstakeEvent) {
      const { returnValues: { undelegatedAt } } = latestUnstakeEvent
      const undelegationComplete = web3Utils.toBN(undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
      setData({
        ...state.data,
        undelegationComplete,
      })
    }
  }, [latestUnstakeEvent])

  return (
    <LoadingOverlay isFetching={isFetching}>
      <section id="pending-undelegation" className="tile">
        <h5>Pending Undelegation</h5>
        <div className="flex pending-undelegation-summary">
          <div className="flex flex-1 flex-column">
            <span className="text-label">TOTAL (KEEP)</span>
            <h2 className="balance flex flex-2">{pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)}`}</h2>
          </div>
          <div className="flex flex-1 flex-column">
            <span className="text-label">UNDELEGATION COMPLETE</span>
            <span className="text-big">{undelegationComplete ? `${undelegationComplete} block` : '-'}</span>
          </div>
          <div className="flex flex-1 flex-column">
            <span className="text-label">UNDELEGATION PERIOD</span>
            <span className="text-big">{undelegationPeriod} blocks</span>
          </div>
        </div>
        <div>
          <PendingUndelegationList pendingUndelegations={pendinUndelegations} />
        </div>
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
