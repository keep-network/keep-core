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
      <h3 className="text-grey-60">Undelegations</h3>
      <SpeechBubbleInfo className="mt-1 mb-1">
        <span className="text-bold">Recover</span>&nbsp;undelegated tokens to return them to your token balance.
      </SpeechBubbleInfo>
      <div className="flex row center">
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
      <ul className="flex column">
        {undelegations && undelegations.map(renderUndelegationItem)}
      </ul>
    </section>
  )
}

const UndelegationItem = React.memo(({ undelegation, successUndelegationCallback }) => {
  return (
    <li className="flex row center space-between text-grey-70" style={{ marginBottom: `0.5rem` }}>
      <div className="flex-1">{undelegation.undelegatedAt}</div>
      <div className="flex-1"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1">{displayAmount(undelegation.amount)} KEEP</div>
      <div className="flex-1">
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
