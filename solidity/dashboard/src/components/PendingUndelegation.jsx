import React, { useEffect, useContext } from 'react'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import {
  displayAmount,
  isSameEthAddress,
  isEmptyObj,
  formatDate,
} from '../utils/general.utils'
import { LoadingOverlay } from './Loadable'
import { Web3Context } from './WithWeb3Context'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { PENDING_STATUS } from '../constants/constants'
import moment from 'moment'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { stakingContract, yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data: {
    pendingUnstakeBalance,
    undelegationCompletedAt,
    undelegationPeriod,
    undelegationStatus,
  } } = state

  useEffect(() => {
    if (!isEmptyObj(latestUnstakeEvent)) {
      const { returnValues: { operator, undelegatedAt } } = latestUnstakeEvent
      if (!isSameEthAddress(yourAddress, operator)) {
        return
      }
      const undelegationCompletedAt = moment.unix(undelegatedAt).add(undelegationPeriod, 'seconds')
      stakingContract.methods.getDelegationInfo(operator).call()
        .then((data) => {
          const { amount } = data
          setData({
            ...state.data,
            undelegationCompletedAt,
            pendingUnstakeBalance: amount,
            undelegationStatus: 'PENDING',
          })
        })
    }
  }, [latestUnstakeEvent.transactionHash])

  const undelegationPeriodRelativeTime = moment().add(undelegationPeriod, 'seconds').fromNow(true)
  const statusText = undelegationStatus === PENDING_STATUS ?
    `${undelegationStatus.toLowerCase()}, ${undelegationCompletedAt.fromNow(true)}` :
    undelegationStatus

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
                text={statusText}
              />
            }
          </div>
          <div className="flex flex-1 column">
            <span className="text-label">completed</span>
            <span className="text-big">{undelegationCompletedAt ? formatDate(undelegationCompletedAt) : '-'}</span>
          </div>
          <div className="flex flex-1 column">
            <span className="text-label">undelegation period</span>
            <span className="text-big">{undelegationPeriodRelativeTime}</span>
          </div>
        </div>
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
