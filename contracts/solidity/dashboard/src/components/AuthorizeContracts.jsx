import React from 'react'
import AddressShortcut from './AddressShortcut'
import { SubmitButton } from './Button'
import { ETHERSCAN_DEFAULT_URL } from '../constants/constants'

const AuthorizeContracts = ({
  contracts,
  onAuthorizeBtn,
  onAuthorizeSuccessCallback,
}) => {
  const renderAuthorizeContractsItem = (item) => (
    <AuthorizeContractsItem
      key={item.contractAddress}
      contract={item}
      onAuthorizeBtn={onAuthorizeBtn}
      onAuthorizeSuccessCallback={onAuthorizeSuccessCallback}
    />
  )

  return (
    <section className="tile">
      <h3 className="text-grey-60">Authorize Contracts</h3>
      <div className="flex row center mt-1">
        <div className="flex-1 text-label">
          contract address
        </div>
        <div className="flex-1 text-label">
          added to the registry
        </div>
        <div className="flex-1 text-label">
            contract details
        </div>
        <div className="flex-1"/>
      </div>
      <ul className="flex column">
        {contracts && contracts.map(renderAuthorizeContractsItem)}
      </ul>
    </section>
  )
}

const AuthorizeContractsItem = ({
  contract,
  onAuthorizeBtn,
  onAuthorizeSuccessCallback,
}) => {
  const authorize = async (onTransactionHashCallback) => {
    await onAuthorizeBtn(contract, onTransactionHashCallback)
  }

  return (
    <li className="flex row center space-between text-grey-70">
      <div className="flex-1">
        <AddressShortcut address={contract.contractAddress} />
      </div>
      <div className="flex-1">
        {contract.blockNumber}
      </div>
      <div className="flex-1">
        <a href={ETHERSCAN_DEFAULT_URL + contract.contractAddress} rel="noopener noreferrer" target="_blank">
        View in Block Explorer
        </a>
      </div>
      <div className="flex-1 flex">
        <SubmitButton
          className="btn btn-primary btn-lg flex-1"
          onSubmitAction={authorize}
          successCallback={onAuthorizeSuccessCallback}
        >
          authorize
        </SubmitButton>
      </div>
    </li>
  )
}

export default AuthorizeContracts
