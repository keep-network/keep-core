import React, { useEffect, useContext } from 'react'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import web3Utils from 'web3-utils'
import { displayAmount, isSameEthAddress, isEmptyObj } from '../utils'
import { LoadingOverlay } from './Loadable'
import { Web3Context } from './WithWeb3Context'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { stakingContract, yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: {
    pendingUnstakeBalance,
    undelegationComplete,
    undelegationPeriod,
    undelegationStatus,
  } } = state

  useEffect(() => {
    if (!isEmptyObj(latestUnstakeEvent)) {
      const { returnValues: { operator, undelegatedAt } } = latestUnstakeEvent
      if (!isSameEthAddress(yourAddress, operator)) {
        return
      }
      const undelegationComplete = web3Utils.toBN(undelegatedAt).add(web3Utils.toBN(undelegationPeriod))
      stakingContract.methods.getDelegationInfo(operator).call()
        .then((data) => {
          const { amount } = data
          setData({
            ...state.data,
            undelegationComplete,
            pendingUnstakeBalance: amount,
            undelegationStatus: 'PENDING',
          })
        })
    }
  }, [latestUnstakeEvent.transactionHash])

  return (
    <LoadingOverlay isFetching={isFetching}>
      <section id="pending-undelegation" className="tile">
        <h3 className="text-grey-60">Token Undelegation</h3>
        <div className="flex pending-undelegation-summary mt-1">
          <div className="flex flex-1 column">
            <span className="text-label">amount</span>
            <h5 className="text-grey-70 flex flex-2">{pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)}`}</h5>
          </div>
          <div className="flex flex-1 column">
            <span className="text-label">undelegation status</span>
            {undelegationStatus &&
              <StatusBadge
                className="self-start"
                status={BADGE_STATUS[undelegationStatus]}
                text={undelegationStatus.toLowerCase()}
              />
            }
          </div>
          <div className="flex flex-1 column">
            <span className="text-label">completed</span>
            <span className="text-big">{undelegationComplete ? `${undelegationComplete} block` : '-'}</span>
          </div>
          <div className="flex flex-1 column">
            <span className="text-label">undelegation period</span>
            <span className="text-big">{undelegationPeriod} blocks</span>
          </div>
        </div>
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
