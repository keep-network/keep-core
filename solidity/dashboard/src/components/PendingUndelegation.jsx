import React, { useEffect, useContext } from 'react'
import { useFetchData } from '../hooks/useFetchData'
import { operatorService } from '../services/token-staking.service'
import web3Utils from 'web3-utils'
import { displayAmount, isSameEthAddress, isEmptyObj } from '../utils/general.utils'
import { LoadingOverlay } from './Loadable'
import { Web3Context } from './WithWeb3Context'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { DataTable, Column } from './DataTable'

const initialData = { pendinUndelegations: [] }

const PendingUndelegation = ({ latestUnstakeEvent }) => {
  const { stakingContract, yourAddress } = useContext(Web3Context)
  const [state, setData] = useFetchData(operatorService.fetchPendingUndelegation, initialData)
  const { isFetching, data } = state
  const {
    undelegationComplete,
    undelegationPeriod,
  } = data

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
        <DataTable data={[data]} itemFieldId="undelegationComplete">
          <Column
            header="amount"
            field="pendingUnstakeBalance"
            renderContent={({ pendingUnstakeBalance }) => (
              pendingUnstakeBalance && `${displayAmount(pendingUnstakeBalance)}`
            )}
          />
          <Column
            header="status"
            field="undelegationStatus"
            renderContent={({ undelegationStatus }) => ( undelegationStatus &&
              <StatusBadge
                className="self-start"
                status={BADGE_STATUS[undelegationStatus]}
                text={undelegationStatus.toLowerCase()}
              />
            )}
          />
          <Column
            header="estimate"
            field="undelegationComplete"
            renderContent={({ undlegationComplete }) =>
              undlegationComplete ? `${undelegationComplete} block` : '-'
            }
          />
          <Column
            header="undelegation period"
            field="undelegationPeriod"
          />
        </DataTable>
      </section>
    </LoadingOverlay>
  )
}

export default PendingUndelegation
