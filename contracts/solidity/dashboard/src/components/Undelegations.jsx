import React from 'react'
import { displayAmount, getAvailableAtBlock } from '../utils'
import AddressShortcut from './AddressShortcut'
import SpeechBubbleInfo from './SpeechBubbleInfo'
import RecoverStakeButton from './RecoverStakeButton'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'


const Undelegations = ({ undelegations, successUndelegationCallback }) => {
  const renderUndelegationItem = (item) =>
    <UndelegationItem
      key={item.operatorAddress}
      undelegation={item}
      successUndelegationCallback={successUndelegationCallback}
    />

  return (
    <section className="tile">
      <h5>Undelegations</h5>
      <SpeechBubbleInfo className="mt-1 mb-1">
        <span className="text-bold">Recover</span>&nbsp;undelegated tokens to return them to your token balance.
      </SpeechBubbleInfo>
      <div className="flex flex-row-center">
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
      <ul className="flex flex-column">
        {undelegations && undelegations.map(renderUndelegationItem)}
      </ul>
    </section>
  )
}

const UndelegationItem = React.memo(({ undelegation, successUndelegationCallback }) => {
  const undelegationStatus = undelegation.canRecoverStake ? COMPLETE_STATUS : PENDING_STATUS

  return (
    <li className="flex flex-row text-darker-grey flex-row-center flex-row-space-between" style={{ marginBottom: `0.8rem` }}>
      <div className="flex-1">{displayAmount(undelegation.amount)} KEEP</div>
      <div className="flex flex-1 flex-column">
        <StatusBadge
          status={BADGE_STATUS[undelegationStatus]}
          className="self-start"
          text={undelegationStatus.toLowerCase()}
        />
        <div className="text-smaller text-darker-grey">
          {getAvailableAtBlock(undelegation.undelegationCompleteAt, undelegationStatus)}
        </div>
      </div>
      <div className="flex-1"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1"><AddressShortcut address={undelegation.authorizerAddress} /></div>
      <div className="flex-1">
        {undelegation.canRecoverStake &&
          <RecoverStakeButton
            successCallback={successUndelegationCallback}
            operatorAddress={undelegation.operatorAddress}
          />
        }
      </div>
    </li>
  )
})

export default Undelegations
