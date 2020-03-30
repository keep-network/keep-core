import React from 'react'
import { displayAmount, formatDate } from '../utils/general.utils'
import AddressShortcut from './AddressShortcut'
import SpeechBubbleInfo from './SpeechBubbleInfo'
import RecoverStakeButton from './RecoverStakeButton'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'


const Undelegations = ({ undelegations }) => {
  const renderUndelegationItem = (item) =>
    <UndelegationItem
      key={item.operatorAddress}
      undelegation={item}
    />

  return (
    <section className="tile">
      <h3 className="text-grey-60">Undelegations</h3>
      <SpeechBubbleInfo className="mt-1 mb-1">
        <span className="text-bold">Recover</span>&nbsp;undelegated tokens to return them to your token balance.
      </SpeechBubbleInfo>
      <div className="flex row center">
        <div className="flex-1 text-label">
          undelegation amount
        </div>
        <div className="flex-1 text-label">
          undelegation status
        </div>
        <div className="flex-1 text-label">
          beneficiary
        </div>
        <div className="flex-1 text-label">
          operator
        </div>
        <div className="flex-1 text-label">
          authorizer
        </div>
        <div className="flex-1" />
      </div>
      <ul className="flex column">
        {undelegations && undelegations.map(renderUndelegationItem)}
      </ul>
    </section>
  )
}

const UndelegationItem = React.memo(({ undelegation }) => {
  const undelegationStatus = undelegation.canRecoverStake ? COMPLETE_STATUS : PENDING_STATUS
  const statusBadgeText = undelegationStatus === PENDING_STATUS ?
    `${undelegationStatus.toLowerCase()}, ${undelegation.undelegationCompleteAt.fromNow(true)}` :
    formatDate(undelegation.undelegationCompleteAt)

  return (
    <li className="flex row center space-between text-grey-70" style={{ marginBottom: `0.5rem` }}>
      <h5 className="flex-1 text-grey-50">{displayAmount(undelegation.amount)} KEEP</h5>
      <div className="flex flex-1 column">
        <StatusBadge
          status={BADGE_STATUS[undelegationStatus]}
          className="self-start"
          text={statusBadgeText}
          onlyIcon={undelegationStatus === COMPLETE_STATUS}
        />
      </div>
      <div className="flex-1"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1"><AddressShortcut address={undelegation.authorizerAddress} /></div>
      <div className="flex-1">
        {undelegation.canRecoverStake &&
          <RecoverStakeButton
            isFromGrant={undelegation.isFromGrant}
            operatorAddress={undelegation.operatorAddress}
          />
        }
      </div>
    </li>
  )
})

export default Undelegations
