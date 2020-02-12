import React from 'react'
import { displayAmount } from '../utils'
import AddressShortcut from './AddressShortcut'

const DelegatedTokensList = ({ delegatedTokens }) => {
  return (
    <section className="tile">
      <h5>Delageted Tokens</h5>
      <div className="flex flex-row-center">
        <div className="flex-1 text-label">
          CREATED AT
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
        {delegatedTokens && delegatedTokens.map(renderDelegatedTokensItem)}
      </ul>
    </section>
  )
}

const DelegatedTokensListItem = ({ undelegation }) => {
  return (
    <li className="flex flex-row flex-row-space-between">
      <div className="flex-1 text-bit">{undelegation.createdAt}</div>
      <div className="flex-1 text-bit"><AddressShortcut address={undelegation.beneficiary} /></div>
      <div className="flex-1 text-bit"><AddressShortcut address={undelegation.operatorAddress} /></div>
      <div className="flex-1 text-bit">{displayAmount(undelegation.amount)} KEEP</div>
    </li>
  )
}

const renderDelegatedTokensItem = (item) => <DelegatedTokensListItem key={item.operatorAddress} undelegation={item} />

export default DelegatedTokensList
