import React, { useEffect, useContext } from 'react'
import PendingUndelegationList from './PendingUndelegationList'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import { formatDate, displayAmount } from '../utils'
import { LoadingOverlay } from './Loadable'
import { Web3Context } from './WithWeb3Context'
import moment from 'moment'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { utils } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: {
    stakeWithdrawalDelayInSec,
    pendingUnstakeBalance,
    undelegatedOn,
    stakeWithdrawalDelay,
    pendinUndelegations,
  } } = state

  useEffect(() => {
    if (latestUnstakeEvent) {
      const { id, returnValues: { createdAt, value } } = latestUnstakeEvent
      const newPendingUndelegation = { eventId: id, createdAt: moment.unix(createdAt), amount: value }
      const updatedPendingUnstakeBalance = utils.toBN(pendingUnstakeBalance).add(utils.toBN(value))
      const updatedUndelegations = [newPendingUndelegation, ...pendinUndelegations]
      const updatedUndelegatedOn = moment.unix(createdAt).add(stakeWithdrawalDelayInSec, 'seconds')

      setData({
        stakeWithdrawalDelayInSec,
        pendingUnstakeBalance: updatedPendingUnstakeBalance,
        undelegatedOn: updatedUndelegatedOn,
        stakeWithdrawalDelay,
        pendinUndelegations: updatedUndelegations,
      })
    }
  }, [latestUnstakeEvent])

  return (
    <LoadingOverlay isFetching={isFetching}>
      <section id="pending-undelegation" className="tile">
        <h5>Pending Undelegation</h5>
        <div className="flex pending-undelegation-summary">
          <h2 className="balance flex flex-2">{pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)} K`}</h2>
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
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
