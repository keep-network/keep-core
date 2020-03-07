import React from 'react'
import AddressShortcut from './AddressShortcut'
import StatusBadge, { BADGE_STATUS } from './StatusBadge'
import { ETHERSCAN_DEFAULT_URL } from '../constants/constants'

const AuthorizationHistory = ({ contracts }) => {
  return (
    <section className="tile">
      <h3 className="text-grey-60">Authorization History</h3>
      <div className="flex row center mt-1">
        <div className="flex-1 text-label">
          contract address
        </div>
        <div className="flex-1 text-label">
          status
        </div>
        <div className="flex-1 text-label">
          contract details
        </div>
      </div>
      <ul className="flex column">
        {contracts && contracts.map(renderAuthorizationHistoryItem)}
      </ul>
    </section>
  )
}

const renderAuthorizationHistoryItem = (item) => (
  <AuthorizationHistoryItem
    key={item.contractAddress}
    contract={item}
  />
)

const AuthorizationHistoryItem = ({ contract }) => {
  return (
    <li className="flex row center space-between text-grey-70">
      <div className="flex-1">
        <AddressShortcut address={contract.contractAddress} />
      </div>
      <div className="flex flex-1">
        <StatusBadge
          className="self-start"
          status={BADGE_STATUS.COMPLETE}
          text="authorized"
        />
      </div>
      <div className="flex-1">
        <a href={ETHERSCAN_DEFAULT_URL + contract.contractAddress} rel="noopener noreferrer" target="_blank">
          View in Block Explorer
        </a>
      </div>
    </li>
  )
}

export default AuthorizationHistory
