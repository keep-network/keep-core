import React, { useEffect, useContext } from 'react'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import web3Utils from 'web3-utils'
import { displayAmount } from '../utils'
import { LoadingOverlay } from './Loadable'
import { Web3Context } from './WithWeb3Context'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { stakingContract } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
  } } = state

  useEffect(() => {
    if (latestUnstakeEvent) {
      const { returnValues: { operator, undelegatedAt } } = latestUnstakeEvent
      const undelegationComplete = web3Utils.toBN(undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
      stakingContract.methods.getUndelegation(operator).call()
        .then((data) => {
          const { amount } = data
          setData({
            ...state.data,
            undelegationComplete,
            pendingUnstakeBalance: amount,
          })
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
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
