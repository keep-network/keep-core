import React from 'react'
import { displayAmount } from '../utils'
import AddressShortcut from './AddressShortcut'
import UndelegateStakeButton from './UndelegateStakeButton'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'

const DelegatedTokensList = ({ delegatedTokens, successDelegationCallback }) => {
  const renderDelegatedTokensItem = (item) =>
    <DelegatedTokensListItem
      key={item.operatorAddress}
      delegation={item}
      successDelegationCallback={successDelegationCallback}
    />

  return (
    <section className="tile">
      <h5>Delageted Tokens</h5>
      <div className="flex flex-row-center mt-1">
        <div className="flex-1 text-label">
          amount
        </div>
        <div className="flex-1 text-label">
          delegation status
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
        <div className="flex-1"/>
      </div>
      <ul className="flex flex-column">
        {delegatedTokens && delegatedTokens.map(renderDelegatedTokensItem)}
      </ul>
    </section>
  )
}

const DelegatedTokensListItem = React.memo(({ delegation, successDelegationCallback }) => {
  const delegationStatus = delegation.isInInitializationPeriod ? 'PENDING' : 'COMPLETE'

  return (
    <li className="flex flex-row text-darker-grey flex-row-center flex-row-space-between" style={{ marginBottom: `0.5rem` }}>
      <div className="flex-1">{displayAmount(delegation.amount)} KEEP</div>
      <div className="flex flex-1">
        <StatusBadge
          status={BADGE_STATUS[delegationStatus]}
          className="self-start"
          text={delegationStatus.toLowerCase()}
        />
      </div>
      <div className="flex-1"><AddressShortcut address={delegation.beneficiary} /></div>
      <div className="flex-1"><AddressShortcut address={delegation.operatorAddress} /></div>
      <div className="flex-1"><AddressShortcut address={delegation.authorizerAddress} /></div>
      <div className="flex-1">
        <UndelegateStakeButton
          isInInitializationPeriod={delegation.isInInitializationPeriod}
          btnClassName="btn btn-sm btn-default"
          operator={delegation.operatorAddress}
          successCallback={successDelegationCallback}
        />
      </div>
    </li>
  )
})

export default DelegatedTokensList
