import React from 'react'
import { displayAmount } from '../utils'
import AddressShortcut from './AddressShortcut'
import SpeechBubbleInfo from './SpeechBubbleInfo'

const Undelegations = ({ undelegations }) => {
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
      </div>
      <ul className="flex flex-column">
        {undelegations && undelegations.map(renderUndelegationItem)}
      </ul>
    </section>
  )
}

const UndelegationItem = ({ undelegation }) => {
  return (
    <li className="flex flex-row flex-row-space-between">
      <div className="flex-1 text-bit">{undelegation.undelegatedAt}</div>
      <div className="flex-1 text-bit"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1 text-bit"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1 text-bit">{displayAmount(undelegation.amount)} KEEP</div>
    </li>
  )
}

const renderUndelegationItem = (item) => <UndelegationItem key={item.operatorAddress} undelegation={item} />

export default Undelegations
