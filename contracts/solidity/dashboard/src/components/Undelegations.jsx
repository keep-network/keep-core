import React from 'react'
import { displayAmount } from '../utils'
import AddressShortcut from './AddressShortcut'
import SpeechBubbleInfo from './SpeechBubbleInfo'
import RecoverStakeButton from './RecoverStakeButton'

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
          UNDELEGATION STARTED
        </div>
        <div className="flex-1 text-label">
          BENEFICIARY ADDRESS
        </div>
        <div className="flex-1 text-label">
          OPERATOR ADDRESS
        </div>
        <div className="flex-1 text-label">
          AMOUNT
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

  return (
    <li className="flex flex-row flex-row-space-between">
      <div className="flex-1 text-big">{undelegation.undelegatedAt}</div>
      <div className="flex-1 text-big"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1 text-big"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1 text-big">{displayAmount(undelegation.amount)} KEEP</div>
      <div className="flex-1 text-big">
        {undelegation.canRecoverStake ?
          <RecoverStakeButton
            successCallback={successUndelegationCallback}
            operatorAddress={undelegation.operatorAddress}
          /> :
          `undelegation will be completed at ${undelegation.undelegationCompleteAt.toString()}`
        }
      </div>
    </li>
  )
})

export default Undelegations
