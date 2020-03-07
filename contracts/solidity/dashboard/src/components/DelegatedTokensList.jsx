import React from 'react'
import { displayAmount, getAvailableAtBlock } from '../utils'
import AddressShortcut from './AddressShortcut'
import UndelegateStakeButton from './UndelegateStakeButton'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { PENDING_STATUS, COMPLETE_STATUS } from '../constants/constants'

const DelegatedTokensList = ({ delegatedTokens, successDelegationCallback }) => {
  const renderDelegatedTokensItem = (item) =>
    <DelegatedTokensListItem
      key={item.operatorAddress}
      delegation={item}
      successDelegationCallback={successDelegationCallback}
    />

  return (
    <section className="tile">
      <h3 className="text-grey-60">Delegations</h3>
      <div className="flex row center mt-1">
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
      <ul className="flex column">
        {delegatedTokens && delegatedTokens.map(renderDelegatedTokensItem)}
      </ul>
    </section>
  )
}

const DelegatedTokensListItem = React.memo(({ delegation, successDelegationCallback }) => {
  const delegationStatus = delegation.isInInitializationPeriod ? PENDING_STATUS : COMPLETE_STATUS

  return (
    <li className="flex row center space-between text-grey-70" style={{ marginBottom: `0.5rem` }}>
      <h5 className="flex-1 text-grey-50">{displayAmount(delegation.amount)} KEEP</h5>
      <div className="flex flex-1 column">
        <StatusBadge
          status={BADGE_STATUS[delegationStatus]}
          className="self-start"
          text={delegationStatus.toLowerCase()}
        />
        <div className="text-smaller text-grey-70">
          {getAvailableAtBlock(delegation.initializationOverAt, delegationStatus)}
        </div>
      </div>
      <div className="flex-1"><AddressShortcut address={delegation.beneficiary} /></div>
      <div className="flex-1"><AddressShortcut address={delegation.operatorAddress} /></div>
      <div className="flex-1"><AddressShortcut address={delegation.authorizerAddress} /></div>
      <div className="flex-1">
        <UndelegateStakeButton
          isInInitializationPeriod={delegation.isInInitializationPeriod}
          btnClassName="btn btn-sm btn-secondary"
          operator={delegation.operatorAddress}
          successCallback={successDelegationCallback}
        />
      </div>
    </li>
  )
})

export default DelegatedTokensList
